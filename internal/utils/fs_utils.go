package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func CopyFile(from, to string) (written int64, err error) {
	dstFile, err := os.Create(to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating the file %s\n", to)
	}
	defer dstFile.Close()

	w := bufio.NewWriter(dstFile)

	srcFile, err := os.Open(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading source file %s\n", from)
	}
	defer srcFile.Close()

	r := bufio.NewReader(srcFile)

	return io.Copy(w, r)
}
