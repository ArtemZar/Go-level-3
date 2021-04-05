package main

import (
	"fmt"
)

func ExampleListDirByReadDir() {
	ListDirByReadDir("../")
	for _, vol := range FindFiles {
		fmt.Printf("File name: %v; File Path: %v; File Size: %d\n", vol.FileName, vol.FilePath, vol.FileSize)
	}
	// Output:
	// File name: file1; File Path: ..//dir1/dir1_2/dir123; File Size: 0
	// File name: file3; File Path: ..//dir1/dir1_2/dir123; File Size: 0
	// File name: file4; File Path: ..//dir1/dir1_2/dir123; File Size: 0
	// File name: file2; File Path: ..//dir1/dir1_2; File Size: 0
	// File name: file4; File Path: ..//dir1; File Size: 0
}

func ExampleFindDubleFiles() {
	FindDubleFiles()
	// Output:
	// Найдены дубликаты файлов:
	// ID: 2; File name: file4; File Path: ..//dir1/dir1_2/dir123; File Size: 0
	// ID: 4; File name: file4; File Path: ..//dir1; File Size: 0
}
