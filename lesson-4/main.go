

package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/ArtemZar/Go-level-3/lesson-4/userQuery"
	"go/parser"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var pathToDataFile string
func init() {
	// задаём и парсим флаги
	flag.StringVar(&pathToDataFile, "path", "data/owid-covid-data.csv", "относительный путь до csv файла")
	flag.Parse()

	// Проверяем обязательное условие
	if pathToDataFile == "" {
		log.Print("no path set by flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	go watchSignals(cancel)

	text := userQuery.Query()
	//fmt.Println(text)

	pes, _ := parser.ParseExpr(text)
	fmt.Printf("%T\n", pes)

file, err := os.Open("data/owid-covid-data.csv")
if err != nil {
panic(err)
}
defer file.Close()

//reader := csv.NewReader(file)
reader := csv.NewReader(bufio.NewReader(file))
//reader.Comma = ','
//reader.FieldsPerRecord = 3
//reader.Comment = '#'

	//for {
	//	line, error := reader.Read()
	//	if error == io.EOF {
	//		break
	//	} else if error != nil {
	//		log.Fatal(error)
	//	}
	//	people = append(people, Person{
	//		Firstname: line[0],
	//		Lastname: line[1],
	//		Address: &Address{
	//			City: line[2],
	//			State: line[3],
	//		},
	//	})
	//}
//var m map [string] []string
	var results [][]string
for i:=0; i<4; i++ {
line, e := reader.Read()
if e != nil {
fmt.Println(e)
break
}
	results = append(results, line)
	ff := Find(results[0], "location")
	if results[i][ff] == "Afghanistan" {
		fmt.Println(results[i])
	}
}





	//fmt.Println(results)
}

// Find возвращает наименьший индекс i,
// при котором x == a[i],
// или len(a), если такого индекса нет.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

// Contains указывает, содержится ли x в a.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func watchSignals(cancel context.CancelFunc) {
	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan,
		syscall.SIGTERM,
		syscall.SIGINT)

	<-osSignalChan
	// если сигнал получен, отменяем контекст работы
	cancel()
}

