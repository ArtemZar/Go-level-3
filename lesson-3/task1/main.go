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
	"fmt"
	"github.com/spf13/afero"
	"hash/crc32"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	//"os/user"
	"runtime"
	"strings"
	"sync"

	logrus "github.com/sirupsen/logrus"
)


type FileList struct {
	FileName string
	FilePath string
	FileSize int64
	FileHash uint32
}

var (
	// флаги
	//del  *bool
	//Path *string

	FindFiles    []FileList // хранит список найденых файлов
	deletedFiles []FileList // хранит только список дубликатов подлежащих удалению
)

// счетчики найденных файлов и каталогов
var countFile, countDir int = 0, 1

// init инициализирует аргументами программы, переданными через командную строку.
//
// не принимает и не возвращает значения
//func init() {
//	del = flag.Bool("d", false, "Accept on del finded duplicate")
//	Path = flag.String("p", "..", "Path to root dir where starting reading files")
//	flag.Parse()
//}



func main() {
	// Функция SetFormatter позволяет установить формат log сообщения - в данном случае JSON
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// создаем корневой логер (далее требуется пробрасывать корневой логер через аргумент функции)

	// Получаем данные текущего пользователя для использования в логах.
	var usr, err = user.Current()
	if err != nil {
	logrus.Fatal("данные о текущем пользователе не получены ", err)
	}
	// преобразовываем относительный путь в абсолютный
	absolutPath, err := filepath.Abs("..")
	if err != nil {
		logrus.Fatal("ошибка преобразования абсолютного пути ", err)
	}

	// Функция Fields добавляет параметры в сообщение
	var standardFields = logrus.Fields{
		"User Name": usr.Username,
		"PID": os.Getpid(),
		"This path": absolutPath,
	}

	// с помощью WithFields создаем логер с параметрами standartFields
	var l = logrus.WithFields(standardFields)

	//

	l.WithFields(logrus.Fields{"func": "main"}).Info("start parsing from ", absolutPath)
	FS := afero.NewOsFs()
	ListDirByReadDir(FS, "..")
	// log.Fields позволяет специфицировать параметры для конкретного сообщения
	l.WithFields(logrus.Fields{"func": "main"}).Info("parsing finish...", "find files:", countFile, " find dir:", countDir)

	// примеры:
	//hlog.WithFields(logrus.Fields{"uid": 100500}).Info("file successfully uploaded")
	//hlog.WithFields(logrus.Fields{"uid": 200512}).Warn("libjpeg: invalid format")
	//hlog.WithFields(logrus.Fields{"uid": 101345}).Error("file corrupted")
	//hlog.Infof("storage space left: %d", 1024)
	FindDubleFiles()
	if  deletedFiles != nil {
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
func ListDirByReadDir(fs afero.Fs, path string) { //, l *logrus.Entry
	l := logrus.WithField("func", "ListDirByReadDir").WithField("PID", os.Getegid())

	// обработка паники
	defer func() {
		if err := recover(); err != nil {
			l.Error("паника при переходе в каталог ", err)

		}
	}()



	mu := sync.Mutex{}
	lst, err := afero.ReadDir(fs, path) // ioutil.ReadDir(path)
	if err != nil {
		l.Error("can't read dir ", path, err)
	}
	for _, val := range lst {
		if !val.IsDir() {
			mu.Lock()
			countFile++
			mu.Unlock()
			hs, err := GetHash(path + "/" + val.Name())
			if err != nil {
			l.Error("can't get hash file: ", val.Name(), "path: ", path, err)
			}
			mu.Lock()
			theFile := FileList{val.Name(), path, val.Size(), hs}
			FindFiles = append(FindFiles, theFile)
			mu.Unlock()
		} else {
			mu.Lock()
			countDir++
			mu.Unlock()
			wg := sync.WaitGroup{}
			wg.Add(1)

			// явный вызов паники
			if x:=(path + "/" + val.Name()); x == "../duble_files" {
				panic(x)
			}

			go func() {
				l.Info("start parsing next dir ", path + "/" + val.Name())
				ListDirByReadDir(fs, path + "/" + val.Name())
				runtime.Gosched()
				wg.Done()
			}()
			wg.Wait()
		}
	}
}

// GetHash функция расчитывает хэш файла.
//
// принимает на вход имя файла (тип string)
//
// возвращает расчитаное значение (тип uint32)
func GetHash(filename string) (uint32, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	h := crc32.NewIEEE()
	h.Write(bs)
	return h.Sum32(), nil
}

// FindDubleFiles функция анализирует срез FindFiles на наличие дубликатов (с линейной сложностью на базе мап)
//
// сравнение производится по полям структуры имя файла (FileName), размер файла (FileSize), хеш файла (FileHash)
//
// не принимает аргументы
//
// возвращаемого значения нет. Итог работы функции вывод найденых дубликатов в стандартный вывод и
//
// формирование списка дубликатов  в срезе deletedFiles для удаления
func FindDubleFiles() {
	type uniqueList struct {
		FileName string
		FileSize int64
		FileHash uint32
	}
	uniqueFile := make(map[uniqueList]string)
	for _, val := range FindFiles {
		keyForMapWithoutPath := uniqueList{val.FileName, val.FileSize, val.FileHash}
		if _, ok := uniqueFile[keyForMapWithoutPath]; !ok {
			uniqueFile[keyForMapWithoutPath] = val.FilePath
			continue
		}
		fmt.Println("Найдены дубликаты файла:")
		fmt.Printf("File name: %v; File Size: %d; File Hash: %d\n", val.FileName, val.FileSize, val.FileHash)
		fmt.Printf("First file path: %v\n", val.FilePath)
		fmt.Printf("Second file path: %v\n", uniqueFile[keyForMapWithoutPath])
		deletedFiles = append(deletedFiles, val)
	}
	if len(deletedFiles)==0 {
		fmt.Println("Дубликаты не найдены")
		return
	}
}

// deletingFiles функция удалят вайлы из операционной системы в соответсвии со списком deletedFiles
//
// не принимает аргументы
//
// возвращаемого значения нет.
func deletingFiles() {
	l := logrus.WithField("func", "deletingFiles").WithField("PID", os.Getegid())
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
			err := os.Remove(vol.FilePath + "/" + vol.FileName)
			if err != nil {
				l.Error("can't remove file: ", vol.FilePath + "/" + vol.FileName, err)
				return
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