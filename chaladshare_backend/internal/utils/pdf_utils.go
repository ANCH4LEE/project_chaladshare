package utils

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

func init() {
	common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
}

func ConvertPDFToImages(pdfPath, outputDir string) ([]string, error) {
	// สร้างโฟลเดอร์
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("สร้างโฟลเดอร์ภาพไม่สำเร็จ: %v", err)
	}

	// เปิด PDF
	f, err := os.Open(pdfPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	var imageFiles []string
	baseName := filepath.Base(pdfPath)
	baseName = baseName[:len(baseName)-4] // ตัด .pdf

	// Render แต่ละหน้า
	device := render.NewImageDevice()

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return nil, err
		}

		img, err := device.Render(page)
		if err != nil {
			return nil, err
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_page_%d.png", baseName, i))

		imgFile, err := os.Create(outputFile)
		if err != nil {
			return nil, err
		}

		err = png.Encode(imgFile, img)
		imgFile.Close()
		if err != nil {
			return nil, err
		}

		imageFiles = append(imageFiles, outputFile)
	}

	return imageFiles, nil
}
