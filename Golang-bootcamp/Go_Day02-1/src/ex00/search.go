package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func SearchFunc(path string, info os.FileInfo, err error) error {
	if os.IsPermission(err) {
		return filepath.SkipDir
	} else if err != nil {
		return err
	}

	if !hasReadPermission(info) {
		return nil
	}

	if info.Mode().IsRegular() && fl.ext != "" {
		fileExtension := filepath.Ext(path)
		if fileExtension == ("." + fl.ext) {
			fmt.Println(path)
		}
	} else {
		if fl.sl && info.Mode()&os.ModeSymlink != 0 {
			s, _ := filepath.EvalSymlinks(path)
			if _, err := os.Stat(s); err == nil {
				fmt.Println(path, "->", s)
			} else {
				fmt.Println(path, "->", "[broken]")
			}
		}

		if fl.f && info.Mode().IsRegular() {
			fmt.Println(path)
		}

		if fl.d && info.IsDir() {
			fmt.Println(path)
		}
	}

	return nil
}

func hasReadPermission(info os.FileInfo) bool {
	// Check if the file or directory has read permission
	return info.Mode().Perm()&0400 != 0
}
