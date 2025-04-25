package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"

	"github.com/gocolly/colly/v2"
)

func WebSource(url string, headers map[string]string, cookies map[string]string) string {
	c := colly.NewCollector()
	source := ""

	if len(headers) > 0 {
		for key, value := range headers {
			c.Headers.Add(key, value)
		}
	}

	if len(cookies) > 0 {
		jar, err := cookiejar.New(nil)
		if err != nil {
			fmt.Println("Error:", err)
		}
		for key, value := range cookies {
			parsedURL, err := neturl.Parse(url)
			if err != nil {
				fmt.Println("Error:", err)
			}
			jar.SetCookies(parsedURL, []*http.Cookie{
				{Name: key, Value: value},
			})
		}
		c.SetCookieJar(jar)
	}

	c.OnResponse(func(r *colly.Response) {
		rawHTML := string(r.Body)
		source = rawHTML
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return source
}

type WebSearchResult struct {
	Position int
	Title    string
	Link     string
	Snippet  string
}

func WebSearch(query string) []WebSearchResult {
	c := colly.NewCollector()
	results := []WebSearchResult{}

	c.OnHTML(".result", func(e *colly.HTMLElement) {
		results = append(results, WebSearchResult{
			Position: e.Index,
			Title:    e.ChildText(".result__title a"),
			Link:     e.ChildAttr(".result__a", "href"),
			Snippet:  e.ChildText(".result__snippet"),
		})
	})

	err := c.Visit("https://html.duckduckgo.com/html?q=" + query)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return results
}
