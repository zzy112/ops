package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"img2pdf/internal/converter"
	"img2pdf/internal/imgutil"
	"img2pdf/internal/model"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectImages opens a file dialog for selecting images
func (a *App) SelectImages() ([]model.FileInfo, error) {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择图片",
		Filters: []runtime.FileFilter{
			{DisplayName: "图片文件", Pattern: "*.jpg;*.jpeg;*.png;*.bmp;*.tiff;*.tif;*.webp;*.gif"},
		},
	})
	if err != nil {
		return nil, err
	}

	var result []model.FileInfo
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		thumb, _ := imgutil.ThumbnailFromFile(f, 200)
		result = append(result, model.FileInfo{
			Path:      f,
			Name:      filepath.Base(f),
			Size:      info.Size(),
			Thumbnail: thumb,
		})
	}
	return result, nil
}

// SelectPDF opens a file dialog for selecting a PDF
func (a *App) SelectPDF() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择PDF文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "PDF文件", Pattern: "*.pdf"},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

// SelectOutputFile opens a save dialog for PDF output
func (a *App) SelectOutputFile() (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "保存PDF",
		DefaultFilename: "output.pdf",
		Filters: []runtime.FileFilter{
			{DisplayName: "PDF文件", Pattern: "*.pdf"},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

// SelectOutputDir opens a directory dialog
func (a *App) SelectOutputDir() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择输出目录",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetPDFPageCount returns the number of pages in a PDF
func (a *App) GetPDFPageCount(path string) (int, error) {
	return converter.GetPDFPageCount(path)
}

// PreviewImage returns a base64 thumbnail of an image
func (a *App) PreviewImage(path string) (string, error) {
	return imgutil.ThumbnailFromFile(path, 400)
}

// PreviewPDFPage renders a PDF page as a base64 thumbnail
func (a *App) PreviewPDFPage(path string, page int) (string, error) {
	img, err := converter.RenderPDFPage(path, page, 150)
	if err != nil {
		return "", err
	}
	return imgutil.GenerateThumbnail(img, 400), nil
}

// ConvertImagesToPDF converts images to a single PDF
func (a *App) ConvertImagesToPDF(images []string, output string, pageSize string) error {
	if len(images) == 0 {
		return fmt.Errorf("请选择至少一张图片")
	}
	if output == "" {
		return fmt.Errorf("请选择输出路径")
	}
	// Ensure .pdf extension
	if !strings.HasSuffix(strings.ToLower(output), ".pdf") {
		output += ".pdf"
	}
	return converter.ImagesToPDF(images, output, pageSize)
}

// ConvertPDFToImages converts PDF pages to images
func (a *App) ConvertPDFToImages(pdfPath string, outputDir string, format string, dpi int, pageRange string) ([]string, error) {
	if pdfPath == "" {
		return nil, fmt.Errorf("请选择PDF文件")
	}
	if outputDir == "" {
		return nil, fmt.Errorf("请选择输出目录")
	}
	if dpi <= 0 {
		dpi = 150
	}
	if format == "" {
		format = "jpg"
	}
	return converter.PDFToImages(pdfPath, outputDir, format, dpi, pageRange)
}

// SelectDocument opens a file dialog for selecting documents
func (a *App) SelectDocument(docType string) (string, error) {
	var filters []runtime.FileFilter
	switch docType {
	case "txt":
		filters = []runtime.FileFilter{{DisplayName: "文本文件", Pattern: "*.txt"}}
	case "html":
		filters = []runtime.FileFilter{{DisplayName: "HTML文件", Pattern: "*.html;*.htm"}}
	case "word":
		filters = []runtime.FileFilter{{DisplayName: "Word文档", Pattern: "*.docx;*.doc"}}
	default:
		filters = []runtime.FileFilter{
			{DisplayName: "所有文档", Pattern: "*.txt;*.html;*.htm;*.docx;*.doc"},
		}
	}
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "选择文档",
		Filters: filters,
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

// SelectSaveFile opens a save dialog with custom filters
func (a *App) SelectSaveFile(defaultName string, filterName string, pattern string) (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "保存文件",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{DisplayName: filterName, Pattern: pattern},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

// ConvertTxtToPDF converts a text file to PDF
func (a *App) ConvertTxtToPDF(inputPath string, outputPath string) error {
	if inputPath == "" {
		return fmt.Errorf("请选择TXT文件")
	}
	if outputPath == "" {
		return fmt.Errorf("请选择输出路径")
	}
	return converter.TxtToPDF(inputPath, outputPath)
}

// ConvertHtmlToPDF converts an HTML file to PDF
func (a *App) ConvertHtmlToPDF(inputPath string, outputPath string) error {
	if inputPath == "" {
		return fmt.Errorf("请选择HTML文件")
	}
	if outputPath == "" {
		return fmt.Errorf("请选择输出路径")
	}
	return converter.HtmlToPDF(inputPath, outputPath)
}

// ConvertDocToPDF converts a Word document to PDF
func (a *App) ConvertDocToPDF(inputPath string, outputPath string) error {
	if inputPath == "" {
		return fmt.Errorf("请选择Word文档")
	}
	if outputPath == "" {
		return fmt.Errorf("请选择输出路径")
	}
	return converter.DocToPDF(inputPath, outputPath)
}

// ConvertPDFToTxt extracts text from PDF
func (a *App) ConvertPDFToTxt(pdfPath string, outputPath string) error {
	if pdfPath == "" {
		return fmt.Errorf("请选择PDF文件")
	}
	if outputPath == "" {
		return fmt.Errorf("请选择输出路径")
	}
	return converter.PDFToTxt(pdfPath, outputPath)
}

// ConvertPDFToHtml converts PDF to HTML
func (a *App) ConvertPDFToHtml(pdfPath string, outputPath string) error {
	if pdfPath == "" {
		return fmt.Errorf("请选择PDF文件")
	}
	if outputPath == "" {
		return fmt.Errorf("请选择输出路径")
	}
	return converter.PDFToHtml(pdfPath, outputPath)
}

// ConvertPDFToDocx converts PDF to DOCX
func (a *App) ConvertPDFToDocx(pdfPath string, outputPath string) error {
	if pdfPath == "" {
		return fmt.Errorf("请选择PDF文件")
	}
	if outputPath == "" {
		return fmt.Errorf("请选择输出路径")
	}
	return converter.PDFToDocx(pdfPath, outputPath)
}

// GetImageInfo returns file info with thumbnail for drag-dropped files
func (a *App) GetImageInfo(path string) (*model.FileInfo, error) {
	if !imgutil.IsSupported(path) {
		return nil, fmt.Errorf("不支持的图片格式: %s", filepath.Ext(path))
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	thumb, _ := imgutil.ThumbnailFromFile(path, 200)
	return &model.FileInfo{
		Path:      path,
		Name:      filepath.Base(path),
		Size:      info.Size(),
		Thumbnail: thumb,
	}, nil
}
