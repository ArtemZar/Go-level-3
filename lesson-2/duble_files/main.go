// Package main implements functions to find and delete duplicate of files in root dir and subdirs
//
// Default value of root dir is "." (dir is from start the program).
//
// Change this param you can use argument -p when starting the program.
//
// Next param -d is accepting on delete finded duplicate of files.
//
// The ListDirByReadDir create file list
//
// ListDirByReadDir(string)
//
// The FindDubleFiles analise file list and  find duplicate of files
//
// FindDubleFiles()
//
// The deletingFiles delete duplicate of files
//
// deletingFiles()

package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

type FileList struct {
	FileName string
	FilePath string
	FileSize int64
	FileHash uint32
}

var (
	// флаги
	del  *bool
	Path *string

	FindFiles    []FileList // хранит список найденых файлов
	deletedFiles []FileList // хранит только список дубликатов подлежащих удалению
)

// init инициализирует аргументами программы, переданными через командную строку.
//
// не принимает и не возвращает значения
func init() {
	del = flag.Bool("d", false, "Accept on del finded duplicate")
	Path = flag.String("p", "../", "Path to root dir where starting reading files")
	flag.Parse()
}

func main() {
	ListDirByReadDir(*Path)
	FindDubleFiles()

	if *del && deletedFiles != nil {
		deletingFiles()
	}
}

// ListDirByReadDir рекурсивная функция парсинга заданного каталога (включая подкаталоги).
//
// Рекурсия вызывается отдельными потоками при перемещении на нижестоящий уровень дерева каталогов (если в каталоге есть подкаталоги)
//
// принимает на вход адрес верхнеуровнего каталога для начала поиска (тип string)
//
// возвращаемого значения нет. Итог работы функции формирование списка файлов в срезе findFiles
func ListDirByReadDir(path string) {
	mu := sync.Mutex{}
	lst, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range lst {
		if !val.IsDir() {
			hs, err := getHash(path + "/" + val.Name())
			if err != nil {
				log.Fatal(err)
			}
			mu.Lock()
			theFile := FileList{val.Name(), path, val.Size(), hs}
			FindFiles = append(FindFiles, theFile)
			mu.Unlock()
		} else {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				ListDirByReadDir(path + "/" + val.Name())
				runtime.Gosched()
				wg.Done()
			}()
			wg.Wait()
		}
	}
}

// getHash функция расчитывает хэш файла.
//
// принимает на вход имя файла (тип string)
//
// возвращает расчитаное значение (тип uint32)
func getHash(filename string) (uint32, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	h := crc32.NewIEEE()
	h.Write(bs)
	return h.Sum32(), nil
}

// FindDubleFiles функция анализирует срез FindFiles на наличие дубликатов
//
// сравнение производится по полям структуры имя файла (FileName), размер файла (FileSize), хеш файла (FileHash)
//
// не принимает аргументы
//
// возвращаемого значения нет. Итог работы функции вывод найденых дубликатов в стандартный вывод и
//
// формирование списка дубликатов  в срезе deletedFiles для удаления
func FindDubleFiles() {
	for ex, vol := range FindFiles {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(ex int, vol FileList) {
			for i := ex + 1; i < len(FindFiles); i++ {

				if vol.FileName == FindFiles[i].FileName &&
					vol.FileSize == FindFiles[i].FileSize &&
					vol.FileHash == FindFiles[i].FileHash {
					fmt.Println("Найдены дубликаты файлов:")
					fmt.Printf("ID: %d; File name: %v; File Path: %v; File Size: %d; File Hash: %d\n", ex, vol.FileName, vol.FilePath, vol.FileSize, vol.FileHash)
					fmt.Printf("ID: %d; File name: %v; File Path: %v; File Size: %d; File Hash: %d\n", i, FindFiles[i].FileName, FindFiles[i].FilePath, FindFiles[i].FileSize, FindFiles[i].FileHash)
					deletedFiles = append(deletedFiles, FindFiles[i])
				}
			}
			wg.Done()
		}(ex, vol)
		wg.Wait()
	}
}

// deletingFiles функция удалят вайлы из операционной системы в соответсвии со списком deletedFiles
//
// не принимает аргументы
//
// возвращаемого значения нет.
func deletingFiles() {
controlQuestion:
	fmt.Println("Вы точно хотите удалить дублирующиеся файлы (y/n)")
	var answer string
	_, err := fmt.Scan(&answer)
	if err != nil {
		fmt.Println("Неверное значение")
		goto controlQuestion
	}
	switch strings.ToLower(answer) {
	case "y":
		for _, vol := range deletedFiles {
			e := os.Remove(vol.FilePath + "/" + vol.FileName)
			if e != nil {
				log.Fatal(e)
			}
			fmt.Printf("Удален файл %v/%v\n", vol.FilePath, vol.FileName)
		}
		deletedFiles = nil
	case "n":
		fmt.Println("удаление файлов отменено")
		break
	default:
		fmt.Println("Неверное значение")
		goto controlQuestion
	}
}