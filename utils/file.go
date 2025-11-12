package utils

import (
	"fmt"
	"io" // WAJIB! buat io.Copy
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// Save file using Fiber's fileHeader
func SaveFiberFile(fileHeader *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	basename := fileHeader.Filename[:len(fileHeader.Filename)-len(ext)]
	filename := fmt.Sprintf("%d-%s%s", time.Now().Unix(), basename, ext)

	// Create uploads folder if not exists
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	// Create destination file
	dst, err := os.Create(filepath.Join("./uploads", filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file data
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, nil
}

// Delete File safely
func DeleteFile(filename string) error {
	if filename == "" {
		return nil
	}

	fullpath := filepath.Join("./uploads", filename)
	if _, err := os.Stat(fullpath); err == nil {
		return os.Remove(fullpath)
	}
	return nil
}
