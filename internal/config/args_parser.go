package config

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type CopyingConfiguration struct {
	MinFileSize          int64
	MaxFileSize          int64
	Verbose              bool
	SrcDirAbsPath        string
	DstDirAbsPath        string
	ExcludedDirPathes    []string
	ExludedFilePathes    []string
	ExcludedCommonPathes []string
}

func (cc *CopyingConfiguration) ParseArgs() bool {
	flag.BoolVar(&cc.Verbose, "verbose", false, "-v=1 | -v=0 | -v | -v=true | -v=F")

	var (
		srcDirRelPath = flag.String("src-directory", "", "--src-directory /media/user/HDD")
		dstDirRelPath = flag.String("dst-directory", "", "--dst-directory /media/user/SSD")

		minFileSizeStrFlag = flag.String("min-file-size", "", "--min-file-size (1000B|1500K|1500M|1500G|1500T)")
		maxFileSizeStrFlag = flag.String("max-file-size", "", "--max-file-size (1000B|1500K|1500M|1500G|1500T)")

		excludedDirPathListFlag = flag.String(
			"exclude-dir-path", "",
			"--exclude-dir-path /some/path/1:/some/path/2:/some/path/3",
		)
		excludedFilePathListFlag = flag.String(
			"exclude-file-path", "",
			"--exclude-file-path file1:file2:file3",
		)
		excludeCommonNamesListFlag = flag.String(
			"exclude-common-names", "",
			"--exclude-common-names name1:name2:name3",
		)
	)

	flag.Parse()

	if *srcDirRelPath == "" || *dstDirRelPath == "" {
		return false
	}

	var err error
	const absError string = "parse args: error while appling Abs func with '%s' arg"

	cc.SrcDirAbsPath, err = filepath.Abs(*srcDirRelPath)
	errorCheck(err, fmt.Sprintf(absError, "srcDirRelPath"))

	cc.DstDirAbsPath, err = filepath.Abs(*dstDirRelPath)
	errorCheck(err, fmt.Sprintf(absError, "dstDirRelPath"))

	cc.MinFileSize = mustParseFileSize(*minFileSizeStrFlag)
	cc.MaxFileSize = mustParseFileSize(*maxFileSizeStrFlag)

	cc.ExcludedDirPathes = ConvertStrToAbsPathList(excludedDirPathListFlag, cc.SrcDirAbsPath)
	cc.ExludedFilePathes = ConvertStrToAbsPathList(excludedFilePathListFlag, cc.SrcDirAbsPath)
	cc.ExcludedCommonPathes = filepath.SplitList(*excludeCommonNamesListFlag)

	success, errors := checkPathes(cc.ExcludedDirPathes, cc.ExludedFilePathes)
	if !success {
		for k, v := range errors {
			fmt.Printf("%s: %v\n", k, v)
		}
		os.Exit(1)
	}

	return true
}

func (cc CopyingConfiguration) String() string {
	return fmt.Sprintf(
		"srcDirAbsPath: %s\n"+
			"dstDirAbsPath: %s\n"+
			"excludedDirPathes: %v\n"+
			"exludedFilePathes: %v\n"+
			"excludedCommonPathes: %v\n"+
			"MinFileSize: %d bytes\n"+
			"MaxFileSize: %d bytes\n"+
			"verbose: %v\n",
		cc.SrcDirAbsPath,
		cc.DstDirAbsPath,
		cc.ExcludedDirPathes,
		cc.ExludedFilePathes,
		cc.ExcludedCommonPathes,
		cc.MinFileSize,
		cc.MaxFileSize,
		cc.Verbose,
	)
}

func ConvertStrToAbsPathList(multiPathStr *string, srcDirBasePath string) []string {
	var pathList []string = filepath.SplitList(*multiPathStr)
	addBasePathIfNeeds(pathList, srcDirBasePath)
	return pathList
}

func addBasePathIfNeeds(pathList []string, srcDirBasePath string) {
	for i := range pathList {
		var err error

		// todo: case when rel path includes source dir name, ex. data_1/folder_1/folder_2
		if strings.HasPrefix(pathList[i], filepath.Base(srcDirBasePath)) {
			pathList[i], err = filepath.Abs(pathList[i])
			if err != nil {
				log.Fatalf("не удалось применить функцию filepath.Abs: %v", err)
			}
		}

		// случай, когда относительный путь указан без доп. указания исходной директории, ex. folder_1/folder_2
		if !filepath.IsAbs(pathList[i]) {
			pathList[i] = filepath.Join(srcDirBasePath, pathList[i])
		} // todo: эта фича - баг! она не позволит исключить директорию с таким же именем как у директории источника!

	}
}

func checkPathes(pathLists ...[]string) (success bool, errors map[string]error) {
	errors = make(map[string]error)
	for _, pathList := range pathLists {
		for i := range pathList {
			if _, err := os.Stat(pathList[i]); err != nil {
				errors[pathList[i]] = err
			}
		}
	}
	if len(errors) == 0 {
		return true, nil
	}
	return false, errors
}

func errorCheck(err error, errMsg string) {
	if err != nil {
		log.Fatalf(errMsg+": %v", err)
	}
}

func mustParseFileSize(value string) int64 {
	if value == "" {
		return 0
	}
	var (
		fileSizeMatchGroup []string = regexp.MustCompile(`(\d+)([BKMGT])`).FindStringSubmatch(value)
		ratio              int64
	)
	switch fileSizeMatchGroup[2] {
	case "B":
		ratio = 1
	case "K":
		ratio = int64(math.Pow(2.0, 10.0))
	case "M":
		ratio = int64(math.Pow(2.0, 20.0))
	case "G":
		ratio = int64(math.Pow(2.0, 30.0))
	case "T":
		ratio = int64(math.Pow(2.0, 40.0))
	}
	paramValue, err := strconv.ParseInt(fileSizeMatchGroup[1], 10, 64)
	if err != nil {
		log.Fatalf("mustParseFileSize: %v", err)
	}
	return paramValue * ratio
}
