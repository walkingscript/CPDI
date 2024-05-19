package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	// simple command line params section
	verboseOutputFlag = flag.Bool("v", false, "set -v before --exclude-dirs and --exclude-files args")
	minFileSizeFlag   = flag.String("min-file-size", "", "ex: rcpd --min-file-size (1000B|1500K|1500M|1500G|1500T)")
	maxFileSizeFlag   = flag.String("max-file-size", "", "ex: rcpd --max-file-size (1000B|1500K|1500M|1500G|1500T)")
	srcFlag           = flag.String("src", "", "-src C:\\Docs")
	dstFlag           = flag.String("dst", "", "-dst F:\\")
	_                 = flag.String("exclude-dirs", "", "--exclude-dirs dir1 dir2 dir3")      // just to suppress panic and give info
	_                 = flag.String("exclude-files", "", "--exclude-files file1 file2 file3") // just to suppress panic and give info

	minFileSizeInBytes int64 // files smaller than this size must not be copied
	maxFileSizeInBytes int64 // files larger than this size must not be copied

	excludedDirs  = make([]string, 0, 10)
	excludedFiles = make([]string, 0, 10)
)

func main() {
	copyDirsAndFilesR(*srcFlag, *dstFlag)
}

func copyDirsAndFilesR(src, dst string) {
	srcBase := filepath.Base(src)
	dstPath := path.Join(dst, srcBase)
	os.MkdirAll(dstPath, fs.ModePerm)
	subEntries, _ := os.ReadDir(src)

	for _, entry := range subEntries {

		if entry.IsDir() {
			absDirPath := path.Join(src, entry.Name())
			if isDirMustBeCopied(absDirPath) {
				copyDirsAndFilesR(absDirPath, dstPath)
			}
			continue
		}

		srcFilePath := filepath.Clean(path.Join(src, entry.Name()))
		dstFilePath := filepath.Clean(path.Join(dst, srcBase, entry.Name()))

		if mbc, code := isFileMustBeCopied(srcFilePath); mbc {
			if _, err := copyFile(srcFilePath, dstFilePath); err == nil {
				// console log
				if *verboseOutputFlag {
					fmt.Printf("copied '%s' -> '%s'\r\n", srcFilePath, dstFilePath)
				}
			} else {
				// console log
				fmt.Fprintf(os.Stderr, "error while coping file '%s' -> '%s': %v\r\n", srcFilePath, dstFilePath, err)

			}
		} else {
			if *verboseOutputFlag {
				switch code {
				case 1:
					fmt.Printf("file '%s' has been ignored because it is excluded by absolute path\r\n", srcFilePath)
				case 2:
					fmt.Printf("file '%s' has been ignored because it is excluded by global path\r\n", srcFilePath)
				case 3:
					fmt.Printf("file '%s' has been ignored because it is excluded by its size\r\n", srcFilePath)
				}
			}
		}
	}
}

// Func isDirMustBeCopied defines by command line args wheather the dir shoould be copied or not
//
// Takes only one arg: - absDirPath: should be clean path! see: filepath.Clean
//
// returns true if the dir should be copied false othewise.
func isDirMustBeCopied(absDirPath string) bool {
	for _, excludedDirPath := range excludedDirs {
		if filepath.IsAbs(excludedDirPath) {
			if _, err := os.Stat(absDirPath); err == nil {
				absDirPath = filepath.Clean(absDirPath)
				if absDirPath == filepath.Clean(excludedDirPath) {

					// console log
					if *verboseOutputFlag {
						fmt.Printf("directory '%s' will be ignored according to request\r\n", absDirPath)
					}

					return false
				}
			}
		} else {
			if strings.Contains(excludedDirPath, "\\") || strings.Contains(excludedDirPath, "/") {
				if filepath.Clean(excludedDirPath) == filepath.Clean(absDirPath) {

					// console log
					if *verboseOutputFlag {
						fmt.Printf("'%s' dir has been ignored because it is excluded by request\r\n", absDirPath)
					}

					return false
				}
			} else {
				if filepath.Base(excludedDirPath) == filepath.Base(absDirPath) {

					// console log
					if *verboseOutputFlag {
						fmt.Printf("'%s' dir has been ignored! all the dirs with name '%s' will be ignored according to request\r\n", absDirPath, excludedDirPath)
					}

					return false
				}
			}

		}
	}
	return true
}

func isFileMustBeCopied(absFilePath string) (result bool, code int) {
	fileInfo, statErr := os.Stat(absFilePath)

	for _, excludedFilePath := range excludedFiles {
		if filepath.IsAbs(excludedFilePath) {
			if statErr == nil {
				if absFilePath == filepath.Clean(excludedFilePath) {
					return false, 1 // absolute path code
				}
			}
		} else {
			if strings.Contains(excludedFilePath, "\\") || strings.Contains(excludedFilePath, "/") {
				if filepath.Clean(excludedFilePath) == filepath.Clean(absFilePath) {

					// console log
					if *verboseOutputFlag {
						fmt.Printf("'%s' dir has been ignored because it is excluded by request\r\n", absFilePath)
					}

					return false, 2 // relative path code
				}
			} else {
				if filepath.Clean(excludedFilePath) == filepath.Base(absFilePath) {

					// console log
					if *verboseOutputFlag {
						fmt.Printf("'%s' dir has been ignored! all the dirs with name '%s' will be ignored according to request\r\n", absFilePath, excludedFilePath)
					}

					return false, 3 // global path code
				}
			}
		}
	}

	if fileInfo.Size() < minFileSizeInBytes || fileInfo.Size() > maxFileSizeInBytes {
		return false, 4
	}

	return true, 0
}

func copyFile(from, to string) (written int64, err error) {
	dstFile, err := os.Create(to)
	defer dstFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creation the file %s\n", to)
	}
	w := bufio.NewWriter(dstFile)

	srcFile, err := os.Open(from)
	defer srcFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading source file %s\n", from)
	}
	r := bufio.NewReader(srcFile)

	return io.Copy(w, r)
}

func init() {
	rcpdParseArgs()
	printStartInfo()
}

func rcpdParseArgs() {
	flag.Parse()
	if *srcFlag == "" || *dstFlag == "" {
		fmt.Fprintln(os.Stderr, "not enougth data to start copying: you should provide data source and destination through -src and -dst params")
		os.Exit(1)
	}
	mustParseComplexArgs()
	minFileSizeInBytes = mustParseFileSize(minFileSizeFlag)
	maxFileSizeInBytes = mustParseFileSize(maxFileSizeFlag)
}

func mustParseFileSize(value *string) int64 {
	if *value == "" {
		return 0
	}
	var (
		fileSizeStr []string = regexp.MustCompile(`(\d+)(\w)`).FindStringSubmatch(*value)
		ratio       int64
	)
	switch fileSizeStr[2] {
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
	paramValue, _ := strconv.ParseInt(fileSizeStr[1], 10, 64)
	maxFileSizeInBytes = paramValue * ratio
	return maxFileSizeInBytes
}

func mustParseComplexArgs() {
	var container *[]string
	for _, arg := range os.Args {
		switch arg {
		case "--exclude-dirs":
			container = &excludedDirs
		case "--exclude-files":
			container = &excludedFiles
		default:
			if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
				container = nil
			}
			if container != nil {
				*container = append(*container, arg)
			}
		}
	}
}

func printStartInfo() {
	fmt.Printf("%s Args %[1]s\r\n", strings.Repeat("=", 30))
	fmt.Println("--min-file-size", *minFileSizeFlag)
	fmt.Println("--max-file-size", *maxFileSizeFlag)
	fmt.Printf(
		"PARSED:\r\nMin-file-size: %d byte\r\nMax-file-size: %d byte\r\nExcluded dirs: %v\r\nExcluded files: %v\r\n",
		minFileSizeInBytes, maxFileSizeInBytes, excludedDirs, excludedFiles,
	)
}
