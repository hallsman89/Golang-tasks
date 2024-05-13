package archiver

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func createArchive(path string, out *os.File) error {
	gw := gzip.NewWriter(out)
	tw := tar.NewWriter(gw)

	err := addData(tw, path)

	if err != nil {
		return err
	}

	errT := tw.Close()
	errG := gw.Close()
	if errT != nil || errG != nil {
		return fmt.Errorf("Error closing archive writers: tar: %v, gzip: %v", errT, errG)
	}

	return nil
}

func addData(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = path
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}
