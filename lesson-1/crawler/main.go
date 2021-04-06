package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	// максимально допустимое число ошибок при парсинге
	errorsLimit = 100000

	// число результатов, которые хотим получить
	resultsLimit = 10000
)

var (
	// адрес в интернете (например, https://en.wikipedia.org/wiki/Lionel_Messi)
	url string

	// насколько глубоко нам надо смотреть (например, 10)
	depthLimit int
)

// Как вы помните, функция инициализации стартует первой
func init() {
	// задаём и парсим флаги
	flag.StringVar(&url, "url", "https://en.wikipedia.org/wiki/Lionel_Messi", "url address")
	flag.IntVar(&depthLimit, "depth", 3, "max depth for run")
	flag.Parse()

	// Проверяем обязательное условие
	if url == "" {
		log.Print("no url set by flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	started := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	/*
	Доработать программу из практической части так,
	чтобы при отправке ей сигнала SIGUSR1 она увеличивала глубину поиска на 10.
	*/
	crawler := newCrawler(depthLimit)
	go watchSignals(crawler, cancel)


	// создаём канал для результатов
	results := make(chan crawlResult)

	// запускаем горутину для чтения из каналов
	done := watchCrawler(ctx, results, errorsLimit, resultsLimit)

	// запуск основной логики
	// внутри есть рекурсивные запуски анализа в других горутинах
	crawler.run(ctx, url, results, 0)

	// ждём завершения работы чтения в своей горутине

	<-done


	log.Println(time.Since(started))

}

// ловим сигналы выключения
func watchSignals(c *crawler, cancel context.CancelFunc) {
	mu := sync.Mutex{}
	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan,
		syscall.SIGTERM,
		syscall.SIGINT)

	sigusr1Chan := make(chan os.Signal)
	signal.Notify(sigusr1Chan, syscall.SIGUSR1)


	for {
		select {
		case sig :=<-osSignalChan:
			log.Printf("!!!got signal %q", sig.String())
			cancel()

		case <-sigusr1Chan:
			log.Println("!!!got sigusr1")
			mu.Lock()
			c.maxDepth = 10
			mu.Unlock()
			continue
		}
		}
	}



func watchCrawler(ctx context.Context, results <-chan crawlResult, maxErrors, maxResults int) chan struct{} {
	readersDone := make(chan struct{})

	go func() {
		defer close(readersDone)
		for {
			select {
			case <-ctx.Done():
				return

			case result := <-results:
				if result.err != nil {
					maxErrors--
					if maxErrors <= 0 {
						log.Println("max errors exceeded")
						return
					}
					continue
				}

				log.Printf("crawling result: %v", result.msg)
				maxResults--
				if maxResults <= 0 {
					log.Println("got max results")
					return
				}
			}
		}
	}()

	return readersDone

}
