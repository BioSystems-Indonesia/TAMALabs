package rest

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ZipDir zips the contents of srcDir into destZip. Returns error if any.
func ZipDir(srcDir, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			if relPath == "." {
				return nil // skip root
			}
			header := &zip.FileHeader{
				Name:     relPath + "/",
				Method:   zip.Deflate,
			}
			_, err := zipWriter.CreateHeader(header)
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}
