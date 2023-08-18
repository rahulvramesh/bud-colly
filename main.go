package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
)

type APIResult struct {
	APIName     string      `json:"api_name"`
	URL         string      `json:"url"`
	Description string      `json:"description"`
	HTMLContent string      `json:"htmlContent"`
	InnerPages  []InnerPage `json:"inner_pages"`
}

type InnerPage struct {
	PageURL     string `json:"page_url"`
	PageContent string `json:"page_content"`
}

func main() {

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("docs.1inch.io"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./coursera_cache"),
	)

	//detailCollector := c.Clone()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		linkURL := e.Request.AbsoluteURL(e.Attr("href"))

		//fmt.Println("Pre Visit", linkURL)

		err := e.Request.Visit(linkURL)
		if err != nil {
			return
		}

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Println("found")
		//fmt.Println(e.DOM.Html())
		fmt.Println(e.DOM.Text())
	})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://docs.1inch.io/docs/1inch-network-overview")

}
