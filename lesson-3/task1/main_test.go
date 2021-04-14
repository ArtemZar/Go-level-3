package main

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert" //require
	"testing"
)

var appFS = afero.NewMemMapFs()


func TestListDirByReadDir(t *testing.T) {
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
		t.Fatal("Срез хранящий найденные вайлы пуст.")
	}

	// Equal Вместо составления if-выражения можем просто проверить, что ожидаемое значение соответствует
	// возвращенному.
	assert.Equal(t, 4, len(FindFiles),"Ожидаемое количество найденных файлов - 4")
}

func TestGetHash(t *testing.T) {
	afero.WriteFile(appFS, "../file1.txt", []byte("file 1"), 0644)
	_, err := GetHash(appFS, "../file1.txt")

	// NoError Позволяет удостовериться, что функция для заданных параметров не вернула ошибки. В отличие от
	// Nil() функции выведет текст ошибки, что удобно при отладке.
assert.NoError(t, err, "Хеш не расчитан")
}

func TestFindDubleFiles(t *testing.T) {
	file1 := FileList{"file1.txt", "../dir1", 300, 123}
	file2 := FileList{"file1.txt", "../dir2", 300, 123}
	file3 := FileList{"file2.txt", "../dir3", 300, 123}
	file4 := FileList{"file3.txt", "../dir3", 300, 123}
	FindFiles = nil
	FindFiles = append(FindFiles, file1, file2, file3, file4)
	FindDubleFiles()
	// Contains Вместо цикла по массиву или ключам хэш-таблицы можно воспользоваться этой функцией для
	// проверки наличия элемента в объекте, по которому можно итерироваться.
assert.Contains(t, deletedFiles, file2, "Дубликат файла (file1.txt) не сохранен в срезе deletedFiles")
}

