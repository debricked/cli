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
	defer zip.Close(zipWriter)

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
	// Unzip a file, or the first file if multiple in zip archive
	r, err := arc.zip.OpenReader(sourcePath)
	if err != nil {

		return err
	}
	defer arc.zip.CloseReader(r)

	for _, file := range r.File {
		fmt.Println(file.Name)
	}

	f := r.File[1]

	outFile, err := arc.fs.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	fmt.Println("file created: ", outFile.Name())
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	fmt.Println("Copying over content...")
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
