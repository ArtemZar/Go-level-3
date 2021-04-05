package main

import (
	"fmt"
	"testing"
)

func TestListDirByReadDir(t *testing.T)  {

		FindFiles = nil
		ListDirByReadDir(*Path)
		fmt.Println(len(FindFiles))
		if len(FindFiles)!=0 {
			t.Fatal("e r r o r")
		}
}
