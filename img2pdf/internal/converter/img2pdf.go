package converter

import (
	"fmt"
	"image"
	"os"

	"img2pdf/internal/imgutil"

	"github.com/signintech/gopdf"
)

func ImagesToPDF(images []string, output string, pageSize string) error {
	pdf := gopdf.GoPdf{}

	for i, imgPath := range images {
		if !imgutil.IsSupported(imgPath) {
			return fmt.Errorf("不支持的图片格式: %s", imgPath)
		}

		img, err := imgutil.LoadImage(imgPath)
		if err != nil {
			return fmt.Errorf("加载第 %d 张图片失败: %w", i+1, err)
		}

		bounds := img.Bounds()
		w := float64(bounds.Dx())
		h := float64(bounds.Dy())

		if pageSize == "A4" {
			if i == 0 {
				pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
			}
			pdf.AddPage()
			pw, ph := gopdf.PageSizeA4.W, gopdf.PageSizeA4.H
			scale := min(pw/w, ph/h)
			sw, sh := w*scale, h*scale
			x := (pw - sw) / 2
			y := (ph - sh) / 2
			if err := addImageToPDF(&pdf, imgPath, x, y, sw, sh); err != nil {
				return fmt.Errorf("写入第 %d 张图片失败: %w", i+1, err)
			}
		} else {
			// original size: 1px = 0.75pt (96dpi)
			ptW := w * 0.75
			ptH := h * 0.75
			if i == 0 {
				pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: ptW, H: ptH}})
			}
			pdf.AddPageWithOption(gopdf.PageOption{PageSize: &gopdf.Rect{W: ptW, H: ptH}})
			if err := addImageToPDF(&pdf, imgPath, 0, 0, ptW, ptH); err != nil {
				return fmt.Errorf("写入第 %d 张图片失败: %w", i+1, err)
			}
		}
	}

	return pdf.WritePdf(output)
}

func addImageToPDF(pdf *gopdf.GoPdf, imgPath string, x, y, w, h float64) error {
	// gopdf needs to re-read the file, so we just pass the path
	return pdf.Image(imgPath, x, y, &gopdf.Rect{W: w, H: h})
}

// GetImageSize returns width and height of an image file
func GetImageSize(path string) (image.Point, error) {
	f, err := os.Open(path)
	if err != nil {
		return image.Point{}, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return image.Point{}, err
	}
	return image.Point{X: cfg.Width, Y: cfg.Height}, nil
}
