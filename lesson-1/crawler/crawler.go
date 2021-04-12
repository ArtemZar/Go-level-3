package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type crawlResult struct {
	err error
	msg string
}

type crawler struct {
	sync.Mutex
	visited  map[string]string
	maxDepth int
}

func newCrawler(maxDepth int) *crawler {
	return &crawler{
		visited:  make(map[string]string),
		maxDepth: maxDepth,
	}
}

// рекурсивно сканируем страницы
func (c *crawler) run(ctx context.Context,  url string, results chan<- crawlResult, depth int) {
	// просто для того, чтобы успевать следить за выводом программы, можно убрать :)
	time.Sleep(2* time.Second)
	mu := sync.Mutex{}
	/*
		общий таймаут на выполнение следующих операций:
		работа парсера, получений ссылок со страницы, формирование заголовка.
	*/
	ctxTimeOut, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// проверяем что контекст исполнения актуален
	select {
	case <-ctx.Done():
		return


	default:
		// проверка глубины
		mu.Lock()
		if depth >= c.maxDepth {
			mu.Unlock()
			ctx.Done()
			return
		}
		mu.Unlock()

		page, err := parse(ctxTimeOut, url)
		if err != nil {
			// ошибку отправляем в канал, а не обрабатываем на месте
			results <- crawlResult{
				err: errors.Wrapf(err, "parse page %s", url),
			}
			return
		}

			title := PageTitle(ctxTimeOut, page)
		links := pageLinks(ctxTimeOut, nil, page)

		// блокировка требуется, т.к. мы модифицируем мапку в несколько горутин
		c.Lock()
		c.visited[url] = title
		c.Unlock()


		// отправляем результат в канал, не обрабатывая на месте
		results <- crawlResult{
			err: nil,
			msg: fmt.Sprintf("%s -> %s\n", url, title),
		}

		// рекурсивно ищем ссылки
		for link := range links {
			// если ссылка не найдена, то запускаем анализ по новой ссылке
			if c.checkVisited(link) {
				continue
			}

			go c.run(ctx, link, results, depth+1)
		}
	}
}

func (c *crawler) checkVisited(url string) bool {
	c.Lock()
	defer c.Unlock()

	_, ok := c.visited[url]
	return ok
}
