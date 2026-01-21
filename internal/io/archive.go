package io

import (
	"encoding/base64"
	"fmt"
	"path"
)

type IArchive interface {
	ZipFile(sourcePath string, targetPath string, zippedName string) error
	UnzipFile(sourcePath string, targetPath string) error
	B64(sourceName string, targetName string) error
	Cleanup(targetName string) error
	IsValidZip(path string) bool
}

type Archive struct {
	workingDirectory string
	fs               IFileSystem
	zip              IZip
}

func NewArchive(workingDirectory string) *Archive {
	return &Archive{
		fs:  FileSystem{},
		zip: Zip{},
	}
}

func NewArchiveWithStructs(workingDirectory string, fs IFileSystem, zip IZip) *Archive {
	return &Archive{
		fs:  fs,
		zip: zip,
	}
}

func (arc *Archive) ZipFile(sourcePath string, targetPath string, zippedName string) error {
	fs := arc.fs
	zip := arc.zip

	sourceContent, err := fs.ReadFile(sourcePath)
	if err != nil {

		return err
	}
	zipFile, err := fs.Create(targetPath)
	if err != nil {

		return err
	}
	defer fs.CloseFile(zipFile)
	zipWriter := zip.NewWriter(zipFile)
	defer zip.CloseWriter(zipWriter) //nolint

	info, err := fs.StatFile(zipFile)
	if err != nil {

		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {

		return err
	}

	header.Name = zippedName
	header.Method = zip.GetDeflate()

	fileWriter, err := zip.CreateHeader(zipWriter, header)
	if err != nil {

		return err
	}

	_, err = fs.WriteToWriter(fileWriter, sourceContent)
	if err != nil {

		return err
	}

	return err
}

func (arc *Archive) UnzipFile(sourcePath string, targetPath string) error {
	r, err := arc.zip.OpenReader(sourcePath)
	if err != nil {

		return err
	}
	defer arc.zip.CloseReader(r) //nolint

	if len(r.File) != 1 {
		return fmt.Errorf("cannot unzip archive which does not contain exactly one file")
	}

	f := r.File[0]
	outFile, err := arc.fs.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close() //nolint

	rc, err := arc.zip.Open(f)
	if err != nil {
		return err
	}
	defer arc.zip.Close(rc) //nolint

	_, err = arc.fs.Copy(outFile, rc)

	return err
}

func (arc *Archive) B64(sourceName string, targetName string) error {
	fs := arc.fs
	fileContent, err := fs.ReadFile(path.Join(arc.workingDirectory, sourceName))
	if err != nil {

		return err
	}

	targetWriter, err := fs.Create(path.Join(arc.workingDirectory, targetName))
	if err != nil {

		return err
	}
	defer fs.CloseFile(targetWriter)

	encodedFile := base64.StdEncoding.EncodeToString(fileContent)
	_, err = fs.WriteToWriter(targetWriter, []byte(encodedFile))
	if err != nil {

		return err
	}

	return err
}

func (arc *Archive) Cleanup(fileName string) error {

	return arc.fs.Remove(path.Join(arc.workingDirectory, fileName))
}

func (arc *Archive) IsValidZip(path string) bool {
	r, err := arc.zip.OpenReader(path)
	if err != nil {
		return false
	}
	defer arc.zip.CloseReader(r)
	return true
}
