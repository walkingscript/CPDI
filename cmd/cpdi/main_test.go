package main

import (
	"fmt"
	"os"
	"os/exec"
	fp "path/filepath"
	"testing"
)

const (
	DIR_TO_COPY_PATH = "data_1"
	DIR_COPY_DST     = "data_2"
)

var (
	foldersToCreate = []string{
		DIR_TO_COPY_PATH,
		fp.Join(DIR_TO_COPY_PATH, "data_1"),
		fp.Join(DIR_TO_COPY_PATH, "do_not_copy"),
		fp.Join(DIR_TO_COPY_PATH, "folder_1"),
		fp.Join(DIR_TO_COPY_PATH, "folder_1", "do_not_copy"),
		fp.Join(DIR_TO_COPY_PATH, "folder_1", "folder_2_excl"),
	}

	filesToCreate = []struct {
		path    string
		content string
	}{
		{path: fp.Join(DIR_TO_COPY_PATH, "file1.txt"), content: "Hello, World!\n"},
		{path: fp.Join(DIR_TO_COPY_PATH, "data_1", "123.txt"), content: "Ohh nooooo!"},
		{path: fp.Join(DIR_TO_COPY_PATH, "do_not_copy", "secret_file.txt"), content: "very big secret\n"},
		{path: fp.Join(DIR_TO_COPY_PATH, "folder_1", "file2.txt"), content: "Another one\n"},
		{path: fp.Join(DIR_TO_COPY_PATH, "folder_1", "file3_excl.txt"), content: ""},
		{path: fp.Join(DIR_TO_COPY_PATH, "folder_1", "do_not_copy", "some_data.txt"), content: "very important data"},
	}

	dirsToCleanUp = [3]string{"bin", DIR_TO_COPY_PATH, DIR_COPY_DST}

	buildAppCommand = exec.Command("go", "build", "-o", "./bin/cpdi", "cpdi/cmd/cpdi")
	runAppCommand   = exec.Command(
		"./bin/cpdi",
		"--src-directory", DIR_TO_COPY_PATH,
		"--dst-directory", DIR_COPY_DST,
		"--min-file-size", "0B",
		"--max-file-size", "1G",
		"--exclude-dir-path", "folder_1/folder_2_excl:data_1",
		"--exclude-file-path", "folder_1/file3_excl.txt",
		"--exclude-common-names", "do_not_copy",
		"--verbose",
	)
)

func cleanUp() {
	for _, path := range dirsToCleanUp {
		os.RemoveAll(path)
	}
}

func createFoldersAndFiles(t *testing.T) {
	for _, folder := range foldersToCreate {
		if err := os.Mkdir(folder, 0777); err != nil {
			t.Fatalf("can't create '%s' folder: %v", folder, err)
		}
	}

	for _, file := range filesToCreate {
		if fd, err := os.Create(file.path); err == nil {
			if _, err := fmt.Fprint(fd, file.content); err != nil {
				t.Fatalf("can't write data to file '%s': %v", file.path, err)
			}
			fd.Close()
		} else {
			t.Fatalf("can't create file '%s': %v", file.path, err)
		}
	}
}

func buildApp(t *testing.T) {
	if err := buildAppCommand.Start(); err != nil {
		t.Fatalf("error during build command run: %v", err)
	}
	if err := buildAppCommand.Wait(); err != nil {
		t.Fatalf("error waiting for build to finish: %v", err)
	}
}

func runApp(t *testing.T) {
	out, err := runAppCommand.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed: %v\nOutput: %s", err, out)
	} else {
		fmt.Print(string(out))
	}
}

type check struct {
	name      string
	path      string
	isDir     bool
	mustExist bool
}

var checkList = []check{
	{name: "new root folder exists", path: DIR_COPY_DST, isDir: true, mustExist: true},
	{name: "dir with the same name as src dir is excluded", path: fp.Join(DIR_COPY_DST, "data_1"), isDir: true, mustExist: false},
	{name: "common dir name is exluded", path: fp.Join(DIR_COPY_DST, "do_not_copy"), isDir: true, mustExist: false},
	{name: "upstream folder is created", path: fp.Join(DIR_COPY_DST, "folder_1"), isDir: true, mustExist: true},
	{name: "downstream common dir is excluded", path: fp.Join(DIR_COPY_DST, "folder_1", "do_not_copy"), isDir: true, mustExist: false},
	{name: "downstream specified path is excluded", path: fp.Join(DIR_COPY_DST, "folder_1", "folder_2_excl"), isDir: true, mustExist: false},
	{name: "common file created", path: fp.Join(DIR_COPY_DST, "file1.txt"), isDir: false, mustExist: true},
	{name: "file in excluded dir is not created", path: fp.Join(DIR_COPY_DST, "data_1", "123.txt"), isDir: false, mustExist: false},
	{name: "downstream file not created in excluded dir", path: fp.Join(DIR_COPY_DST, "do_not_copy", "secret_file.txt"), isDir: false, mustExist: false},
	{name: "common file created in subdir", path: fp.Join(DIR_COPY_DST, "folder_1", "file2.txt"), isDir: false, mustExist: true},
	{name: "excluded file is not created in subdir", path: fp.Join(DIR_COPY_DST, "folder_1", "file3_excl.txt"), isDir: false, mustExist: false},
	{name: "file is not created in excluded common subdir", path: fp.Join(DIR_COPY_DST, "folder_1", "do_not_copy", "some_data.txt"), isDir: false, mustExist: false},
}

func TestMain(t *testing.T) {
	cleanUp()
	defer cleanUp()

	buildApp(t)
	createFoldersAndFiles(t)
	runApp(t)

	for _, check := range checkList {
		t.Run(check.name, func(t *testing.T) {
			if stat, err := os.Stat(check.path); err != nil {
				if check.mustExist { // were not able to get stat info and path must not exist
					t.Errorf("path '%s' must exist (error: %v)", check.path, err)
				}
			} else {
				if isDir := stat.IsDir(); isDir == check.isDir {
					if !check.mustExist {
						t.Errorf("path '%s' must NOT exist (error: %v)", check.path, err)
					}
				} else {
					t.Errorf("object type (dir/file) not matched for path '%s'", check.path)
				}
			}
		})
	}

}
