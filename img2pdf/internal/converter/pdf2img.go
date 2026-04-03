package converter

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gen2brain/go-fitz"
)

func PDFToImages(pdfPath string, outputDir string, format string, dpi int, pageRange string) ([]string, error) {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("打开PDF失败: %w", err)
	}
	defer doc.Close()

	totalPages := doc.NumPage()
	pages, err := parsePageRange(pageRange, totalPages)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建输出目录失败: %w", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))
	var outputFiles []string

	for _, pageNum := range pages {
		img, err := doc.ImageDPI(pageNum, float64(dpi))
		if err != nil {
			return outputFiles, fmt.Errorf("渲染第 %d 页失败: %w", pageNum+1, err)
		}

		ext := format
		if ext == "jpg" {
			ext = "jpg"
		}
		outPath := filepath.Join(outputDir, fmt.Sprintf("%s_page_%d.%s", baseName, pageNum+1, ext))

		f, err := os.Create(outPath)
		if err != nil {
			return outputFiles, fmt.Errorf("创建文件失败: %w", err)
		}

		switch format {
		case "jpg", "jpeg":
			err = jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
		case "png":
			err = png.Encode(f, img)
		default:
			err = jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
		}
		f.Close()
		if err != nil {
			return outputFiles, fmt.Errorf("保存第 %d 页失败: %w", pageNum+1, err)
		}
		outputFiles = append(outputFiles, outPath)
	}

	return outputFiles, nil
}

func GetPDFPageCount(pdfPath string) (int, error) {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return 0, fmt.Errorf("打开PDF失败: %w", err)
	}
	defer doc.Close()
	return doc.NumPage(), nil
}

func RenderPDFPage(pdfPath string, page int, dpi float64) (image.Image, error) {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("打开PDF失败: %w", err)
	}
	defer doc.Close()

	if page < 0 || page >= doc.NumPage() {
		return nil, fmt.Errorf("页码超出范围: %d", page+1)
	}

	img, err := doc.ImageDPI(page, dpi)
	if err != nil {
		return nil, fmt.Errorf("渲染页面失败: %w", err)
	}
	return img, nil
}

// parsePageRange parses "all", "1-5", "1,3,5", "1-3,5,7-9"
func parsePageRange(rangeStr string, total int) ([]int, error) {
	if rangeStr == "" || rangeStr == "all" {
		pages := make([]int, total)
		for i := range pages {
			pages[i] = i
		}
		return pages, nil
	}

	var pages []int
	parts := strings.Split(rangeStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
			if err != nil {
				return nil, fmt.Errorf("无效的页码范围: %s", part)
			}
			end, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err != nil {
				return nil, fmt.Errorf("无效的页码范围: %s", part)
			}
			if start < 1 || end > total || start > end {
				return nil, fmt.Errorf("页码范围超出: %s (共 %d 页)", part, total)
			}
			for i := start; i <= end; i++ {
				pages = append(pages, i-1) // 0-indexed
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("无效的页码: %s", part)
			}
			if p < 1 || p > total {
				return nil, fmt.Errorf("页码超出范围: %d (共 %d 页)", p, total)
			}
			pages = append(pages, p-1)
		}
	}
	return pages, nil
}
