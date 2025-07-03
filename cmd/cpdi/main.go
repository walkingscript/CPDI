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
	os.Mkdir(params.DstDirAbsPath, 0777)
	CopyDir(params.SrcDirAbsPath)
}

func init() {
	params.ParseArgs()
	if params.Verbose {
		fmt.Print(params)
	}
}

func CopyDir(absSrcPath string) {
	err := os.Chdir(absSrcPath)
	if err != nil {
		log.Fatalf("не получилось открыть директорию: %v", err)
	}
	defer os.Chdir("..")

	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("ошибка чтения директории: %v", err)
	}

	for _, entry := range entries {

		if slices.Contains(params.ExcludedCommonPathes, entry.Name()) {
			continue
		}

		if entry.IsDir() {

			dirAbsPath, err := filepath.Abs(entry.Name())
			if err != nil {
				log.Fatalf("ошибка копирования директории: не удалось применить filepath.Abs: %v", err)
			}
			if slices.Contains(params.ExcludedDirPathes, dirAbsPath) {
				continue
			}

			relPath, err := filepath.Rel(params.SrcDirAbsPath, dirAbsPath)
			if err != nil {
				log.Fatalf("не удалось получить относительные путь: %v", err)
			}
			err = os.Mkdir(filepath.Join(params.DstDirAbsPath, relPath), 0777)
			if err != nil {
				log.Fatalf("ошибка создания директории %s: %v", dirAbsPath, err)
			}
			CopyDir(dirAbsPath)
			continue
		}

		// для файлов

		absFilePath, err := filepath.Abs(entry.Name())
		if err != nil {
			log.Fatalf("не получилось получить абсолютный путь для файла: %v", err)
		}

		if slices.Contains(params.ExludedFilePathes, absFilePath) {
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			log.Fatalf("ошибка получения информации о файле: %v", err)
		}

		if fileInfo.Size() < params.MinFileSize || fileInfo.Size() > params.MaxFileSize {
			continue
		}

		relPath, err := filepath.Rel(params.SrcDirAbsPath, absFilePath)
		if err != nil {
			log.Fatalf("не удалось получить относительные путь: %v", err)
		}
		_, err = utils.CopyFile(absFilePath, filepath.Join(params.DstDirAbsPath, relPath))
		if err != nil {
			log.Fatalf("не удалось скопировать файл: %v", err)
		}
	}
}
