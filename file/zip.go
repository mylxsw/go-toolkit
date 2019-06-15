package file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// ZipFile is a file wrapper contains filename and path
type ZipFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// CreateZipArchiveFile create a zip archive file from files
func CreateZipArchiveFile(saveAs string, files []ZipFile, ignoreError bool) (err error) {
	defer func() {
		if err != nil {
			os.Remove(saveAs)
		}
	}()

	saveAsFile, err := os.Create(saveAs)
	if err != nil {
		return fmt.Errorf("can not create destination zip archive file: %s", err.Error())
	}

	defer saveAsFile.Close()

	_, err = CreateZipArchive(saveAsFile, files, ignoreError)
	return
}

// CreateZipArchiveFileWithIgnored create a zip archive file from files, and return all ignored files
func CreateZipArchiveFileWithIgnored(saveAs string, files []ZipFile, ignoreError bool) (ignoredFiles []ZipFile, err error) {
	defer func() {
		if err != nil {
			os.Remove(saveAs)
		}
	}()

	saveAsFile, err := os.Create(saveAs)
	if err != nil {
		return ignoredFiles, fmt.Errorf("can not create destination zip archive file: %s", err.Error())
	}

	defer saveAsFile.Close()

	return CreateZipArchive(saveAsFile, files, ignoreError)
}

// CreateZipArchive create a zip archive
func CreateZipArchive(w io.Writer, files []ZipFile, ignoreError bool) ([]ZipFile, error) {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	var ignoredFiles = make([]ZipFile, 0)
	for _, file := range files {
		if err := addFileToZipArchive(zipWriter, file); err != nil {
			ignoredFiles = append(ignoredFiles, file)

			if !ignoreError {
				return ignoredFiles, err
			}
		}
	}

	return ignoredFiles, nil
}

// addFileToZipArchive add a file to zip archive
func addFileToZipArchive(zipWriter *zip.Writer, file ZipFile) error {

	zipfile, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("can not open %s: %s", file.Path, err.Error())
	}

	defer zipfile.Close()

	info, err := zipfile.Stat()
	if err != nil {
		return fmt.Errorf("can not get file state for %s: %s", file.Path, err.Error())
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("create zip file header failed for %s: %s", file.Path, err.Error())
	}

	header.Name = file.Name
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create header failed for %s: %s", file.Path, err.Error())
	}

	_, err = io.Copy(writer, zipfile)
	if err != nil {
		return fmt.Errorf("copy file content failed for %s: %s", file.Path, err.Error())
	}

	return nil
}
