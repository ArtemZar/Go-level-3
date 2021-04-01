package main

import (
	"context"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

// парсим страницу
func parse(ctx context.Context, url string) (*html.Node, error) {
		// что здесь должно быть вместо http.Get? :)
	c1:=make(chan bool)
	defer func() {c1<-true}()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-c1:
			break
		}
	}()

		r, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("can't get page")
		}
		b, err := html.Parse(r.Body)
		if err != nil {
			return nil, fmt.Errorf("can't parse page")
		}

		return b, err
}

// ищем заголовок на странице
func PageTitle(ctx context.Context, n *html.Node) string {
	c1:=make(chan bool)
	defer func() {c1<-true}()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-c1:
			break
		}
	}()
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = PageTitle(ctx, c)
		if title != "" {
			break
		}
	}
	return title
}

// ищем все ссылки на страницы. Используем мапку чтобы избежать дубликатов
func pageLinks(ctx context.Context, links map[string]struct{}, n *html.Node) map[string]struct{} {
	c1:=make(chan bool)
	defer func() {c1<-true}()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
		case <-c1:
			break
		}
	}()
	if links == nil {
		links = make(map[string]struct{})
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			// костылик для простоты
			if _, ok := links[a.Val]; !ok && len(a.Val) > 2 && a.Val[:2] == "//" {
				links["http://"+a.Val[2:]] = struct{}{}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = pageLinks(ctx, links, c)
	}
	return links
}
