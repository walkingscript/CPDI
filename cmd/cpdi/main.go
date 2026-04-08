// File Copying System with filtering files and directories by names and by size
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"cpdi/internal/config"
	"cpdi/internal/utils"
)

var params config.CopyingConfiguration

func main() {
	if !params.ParseArgs() {
		fmt.Print(config.HelpString)
		os.Exit(0)
	}
	if params.Verbose {
		fmt.Print(params)
	}
	os.Mkdir(params.DstDirAbsPath, 0777)
	CopyDir(params.SrcDirAbsPath)
}

func CopyDir(absSrcPath string) {
	err := os.Chdir(absSrcPath)
	if err != nil {
		log.Fatalf("error while opening dir: %v", err)
	}
	defer os.Chdir("..")

	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("error while reading dir: %v", err)
	}

	for _, entry := range entries {

		if slices.Contains(params.ExcludedCommonPathes, entry.Name()) {
			continue
		}

		if entry.IsDir() {

			dirAbsPath, err := filepath.Abs(entry.Name())
			if err != nil {
				log.Fatalf("dir copy error: filepath.Abs apply error: %v", err)
			}
			if slices.Contains(params.ExcludedDirPathes, dirAbsPath) {
				continue
			}

			relPath, err := filepath.Rel(params.SrcDirAbsPath, dirAbsPath)
			if err != nil {
				log.Fatalf("getting REL path error: %v", err)
			}
			err = os.Mkdir(filepath.Join(params.DstDirAbsPath, relPath), 0777)
			if err != nil {
				log.Fatalf("error while making dir '%s': %v", dirAbsPath, err)
			}
			CopyDir(dirAbsPath)
			continue
		}

		// для файлов

		absFilePath, err := filepath.Abs(entry.Name())
		if err != nil {
			log.Fatalf("error getting ABS path: %v", err)
		}

		if slices.Contains(params.ExludedFilePathes, absFilePath) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			log.Fatalf("error getting file stats: %v", err)
		}

		if fileInfo.Size() < params.MinFileSize || fileInfo.Size() > params.MaxFileSize {
			continue
		}

		relPath, err := filepath.Rel(params.SrcDirAbsPath, absFilePath)
		if err != nil {
			log.Fatalf("can't get REL path: %v", err)
		}
		dst := filepath.Join(params.DstDirAbsPath, relPath)
		_, err = utils.CopyFile(absFilePath, dst)
		if err != nil {
			log.Fatalf("unable to copy the file '%s' -> '%s': %v", absFilePath, dst, err)
		}
	}
}
