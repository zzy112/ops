package converter

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/signintech/gopdf"
)

// TxtToPDF converts a text file to PDF with Chinese font support
func TxtToPDF(inputPath string, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// Try to load Chinese font (TTF only, TTC not supported by gopdf)
	fontPath := `C:\Windows\Fonts\simhei.ttf`
	if _, err := os.Stat(fontPath); err != nil {
		fontPath = `C:\Windows\Fonts\simkai.ttf`
		if _, err := os.Stat(fontPath); err != nil {
			return fmt.Errorf("未找到中文字体文件（simhei.ttf 或 simkai.ttf）")
		}
	}

	if err := pdf.AddTTFFont("chinese", fontPath); err != nil {
		return fmt.Errorf("加载字体失败: %w", err)
	}
	if err := pdf.SetFont("chinese", "", 11); err != nil {
		return fmt.Errorf("设置字体失败: %w", err)
	}

	pageW := gopdf.PageSizeA4.W
	pageH := gopdf.PageSizeA4.H
	marginX := 40.0
	marginY := 40.0
	lineHeight := 18.0
	maxWidth := pageW - marginX*2
	curY := marginY

	pdf.AddPage()

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			curY += lineHeight
			if curY+lineHeight > pageH-marginY {
				pdf.AddPage()
				curY = marginY
			}
			continue
		}

		// Word-wrap long lines
		for len(line) > 0 {
			if curY+lineHeight > pageH-marginY {
				pdf.AddPage()
				curY = marginY
			}

			// Find how many chars fit in maxWidth
			fitted := fitText(&pdf, line, maxWidth)
			if fitted <= 0 {
				fitted = 1
			}

			pdf.SetX(marginX)
			pdf.SetY(curY)
			pdf.Cell(nil, line[:fitted])
			curY += lineHeight
			line = line[fitted:]
		}
	}

	return pdf.WritePdf(outputPath)
}

// fitText returns how many bytes of text fit within maxWidth
func fitText(pdf *gopdf.GoPdf, text string, maxWidth float64) int {
	runes := []rune(text)
	for i := len(runes); i > 0; i-- {
		w, _ := pdf.MeasureTextWidth(string(runes[:i]))
		if w <= maxWidth {
			return len(string(runes[:i]))
		}
	}
	return len(string(runes[:1]))
}

// HtmlToPDF converts HTML file to PDF using Chrome/Edge headless
func HtmlToPDF(inputPath string, outputPath string) error {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("获取文件路径失败: %w", err)
	}

	fileURL := "file:///" + strings.ReplaceAll(absInput, `\`, `/`)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 60*time.Second)
	defer timeoutCancel()

	var pdfBuf []byte
	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	); err != nil {
		return fmt.Errorf("HTML转PDF失败: %w", err)
	}

	if err := os.WriteFile(outputPath, pdfBuf, 0644); err != nil {
		return fmt.Errorf("保存PDF失败: %w", err)
	}
	return nil
}

// DocToPDF converts Word document to PDF using LibreOffice
func DocToPDF(inputPath string, outputPath string) error {
	// Find LibreOffice
	sofficePath := findLibreOffice()
	if sofficePath == "" {
		return fmt.Errorf("未找到 LibreOffice，请先安装 LibreOffice 以支持 Word 转 PDF")
	}

	outputDir := filepath.Dir(outputPath)
	cmd := exec.Command(sofficePath,
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outputDir,
		inputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("LibreOffice 转换失败: %s\n%s", err, string(output))
	}

	// LibreOffice outputs with the same base name
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	generatedPDF := filepath.Join(outputDir, baseName+".pdf")

	// Rename if output path differs
	if generatedPDF != outputPath {
		if err := os.Rename(generatedPDF, outputPath); err != nil {
			return fmt.Errorf("移动输出文件失败: %w", err)
		}
	}

	return nil
}

func findLibreOffice() string {
	// 1. Check bundled portable LibreOffice next to executable
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)

		// Windows: exe_dir/libreoffice/program/soffice.exe
		bundled := filepath.Join(exeDir, "libreoffice", "program", "soffice.exe")
		if _, err := os.Stat(bundled); err == nil {
			return bundled
		}
		// Windows: exe_dir/libreoffice/program/soffice (no ext)
		bundled = filepath.Join(exeDir, "libreoffice", "program", "soffice")
		if _, err := os.Stat(bundled); err == nil {
			return bundled
		}

		// macOS .app bundle: Contents/MacOS/../Resources/libreoffice
		bundled = filepath.Join(exeDir, "..", "Resources", "libreoffice", "program", "soffice")
		if _, err := os.Stat(bundled); err == nil {
			return bundled
		}
	}

	// 2. System-installed LibreOffice
	candidates := []string{
		`C:\Program Files\LibreOffice\program\soffice.exe`,
		`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		`/Applications/LibreOffice.app/Contents/MacOS/soffice`,
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// 3. Try PATH
	if p, err := exec.LookPath("soffice"); err == nil {
		return p
	}
	return ""
}
