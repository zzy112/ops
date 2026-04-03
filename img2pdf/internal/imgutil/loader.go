package imgutil

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

var SupportedExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".bmp": true, ".tiff": true, ".tif": true,
	".webp": true, ".gif": true,
}

func IsSupported(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return SupportedExts[ext]
}

func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}
	return img, nil
}
