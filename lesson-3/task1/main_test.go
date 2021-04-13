package main

import (
	"github.com/spf13/afero"
	"testing"
)

func TestListDirByReadDir(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("../testDir", 0755)
	appFS.MkdirAll("../testDir/1", 0755)
	appFS.MkdirAll("../testDir/1/1", 0755)
	appFS.MkdirAll("../testDir/1/2", 0755)
	appFS.MkdirAll("../testDir/2", 0755)
	appFS.MkdirAll("../testDir/2/1", 0755)
	appFS.MkdirAll("../testDir/2/2", 0755)
	appFS.MkdirAll("../testDir/3", 0755)
	afero.WriteFile(appFS, "../testDir/1/1/file1.txt", []byte("file 1"), 0644)
	afero.WriteFile(appFS, "../testDir/2/file2.txt", []byte("file 2"), 0644)
	afero.WriteFile(appFS, "../testDir/2/2/file1.txt", []byte("file 1"), 0644)
	afero.WriteFile(appFS, "../testDir/3/file3.txt", []byte("file 3"), 0644)
	ListDirByReadDir(appFS, "..")
	if len(FindFiles)==0 {
		t.Fatal("e r r o r")
	}

}

