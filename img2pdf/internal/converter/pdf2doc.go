package converter

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
)

// PDFToTxt extracts text from PDF and saves as TXT
func PDFToTxt(pdfPath string, outputPath string) error {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("打开PDF失败: %w", err)
	}
	defer doc.Close()

	var sb strings.Builder
	for i := 0; i < doc.NumPage(); i++ {
		text, err := doc.Text(i)
		if err != nil {
			return fmt.Errorf("提取第 %d 页文本失败: %w", i+1, err)
		}
		if i > 0 {
			sb.WriteString("\n--- 第 " + fmt.Sprintf("%d", i+1) + " 页 ---\n\n")
		}
		sb.WriteString(text)
	}

	if err := os.WriteFile(outputPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}
	return nil
}

// PDFToHtml extracts text from PDF and generates a simple HTML file
func PDFToHtml(pdfPath string, outputPath string) error {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("打开PDF失败: %w", err)
	}
	defer doc.Close()

	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>` + escapeHTML(filepath.Base(pdfPath)) + `</title>
<style>
  body { font-family: "Microsoft YaHei", sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; line-height: 1.8; color: #333; }
  .page { margin-bottom: 30px; padding-bottom: 20px; border-bottom: 1px solid #ddd; }
  .page-header { color: #888; font-size: 12px; margin-bottom: 10px; }
  pre { white-space: pre-wrap; word-wrap: break-word; }
</style>
</head>
<body>
`)

	for i := 0; i < doc.NumPage(); i++ {
		text, err := doc.Text(i)
		if err != nil {
			return fmt.Errorf("提取第 %d 页文本失败: %w", i+1, err)
		}
		sb.WriteString(fmt.Sprintf(`<div class="page">
<div class="page-header">第 %d 页</div>
<pre>%s</pre>
</div>
`, i+1, escapeHTML(text)))
	}

	sb.WriteString("</body>\n</html>\n")

	if err := os.WriteFile(outputPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}
	return nil
}

// PDFToDocx extracts text from PDF and creates a DOCX file
func PDFToDocx(pdfPath string, outputPath string) error {
	pdfDoc, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("打开PDF失败: %w", err)
	}
	defer pdfDoc.Close()

	var paragraphs []string
	for i := 0; i < pdfDoc.NumPage(); i++ {
		text, err := pdfDoc.Text(i)
		if err != nil {
			return fmt.Errorf("提取第 %d 页文本失败: %w", i+1, err)
		}
		if i > 0 {
			paragraphs = append(paragraphs, fmt.Sprintf("--- 第 %d 页 ---", i+1))
		}
		lines := strings.Split(text, "\n")
		paragraphs = append(paragraphs, lines...)
	}

	return writeDocx(outputPath, paragraphs)
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func writeDocx(outputPath string, paragraphs []string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("保存DOCX失败: %w", err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// [Content_Types].xml
	writeZipFile(w, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`)

	// _rels/.rels
	writeZipFile(w, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`)

	// word/_rels/document.xml.rels
	writeZipFile(w, "word/_rels/document.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`)

	// word/styles.xml
	writeZipFile(w, "word/styles.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:docDefaults>
    <w:rPrDefault>
      <w:rPr>
        <w:rFonts w:ascii="Calibri" w:eastAsia="SimSun" w:hAnsi="Calibri" w:cs="Times New Roman"/>
        <w:sz w:val="24"/>
        <w:szCs w:val="24"/>
        <w:lang w:val="en-US" w:eastAsia="zh-CN"/>
      </w:rPr>
    </w:rPrDefault>
  </w:docDefaults>
  <w:style w:type="paragraph" w:default="1" w:styleId="Normal">
    <w:name w:val="Normal"/>
  </w:style>
</w:styles>`)

	// word/document.xml
	var body strings.Builder
	body.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>`)
	for _, p := range paragraphs {
		body.WriteString(`<w:p><w:r><w:t xml:space="preserve">`)
		body.WriteString(escapeXML(p))
		body.WriteString(`</w:t></w:r></w:p>`)
	}
	body.WriteString(`</w:body></w:document>`)
	writeZipFile(w, "word/document.xml", body.String())

	return nil
}

func writeZipFile(w *zip.Writer, name, content string) {
	f, _ := w.Create(name)
	f.Write([]byte(content))
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}
