package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/k3a/html2text"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
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

func generateRandomFileName() string {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63()
	return fmt.Sprintf("randomFile%d.html", randNum)
}

func fetchHTMLContent(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching page content: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	bodyString := string(bodyBytes)

	//fmt.Print(bodyString)
	plain := html2text.HTML2Text(bodyString)
	//
	fmt.Println(plain)

	fileName := generateRandomFileName()
	err = ioutil.WriteFile(fileName, []byte(plain), 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}

	fmt.Printf("Saved to file: %s\n", fileName)

	//node, err := html.Parse(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//textContent := extractText(node)
	//fmt.Println(textContent)

}
func processPage(apiURL string) *APIResult {

	fmt.Println("Got URL", apiURL)
	c := colly.NewCollector()

	var result APIResult

	c.OnHTML("body", func(e *colly.HTMLElement) {

		fmt.Print(e.Text)

		result.APIName = "api_name"
		result.URL = apiURL
		result.Description = "description"
		result.HTMLContent = e.Text
	})

	if result.URL != "" {
		fmt.Println(result)
		return &result
	}
	return nil

}

func main() {

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("18f.github.io"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./coursera_cache"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		linkURL := e.Request.AbsoluteURL(e.Attr("href"))
		//parsedURL, err := url.Parse(linkURL)

		//fmt.Println(parsedURL)

		//if err != nil {
		//	log.Printf("Error parsing URL: %v", err)
		//	return
		//}

		err := e.Request.Visit(linkURL)
		if err != nil {
			return
		}

		//processPage(parsedURL.String())
		//if parsedURL.Host == baseDomain {
		//	// Recursively process the linked page
		//	pageResult := processPage(linkURL)
		//	if pageResult != nil {
		//		results = append(results, *pageResult)
		//	}
		//}
	})

	// Create another collector to scrape course details
	//detailCollector := c.Clone()

	//courses := make([]Course, 0, 200)

	// On every <a> element which has "href" attribute call callback
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	// If attribute class is this long string return from callback
	//	// As this a is irrelevant
	//	if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
	//		return
	//	}
	//	link := e.Attr("href")
	//	// If link start with browse or includes either signup or login return from callback
	//	if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
	//		return
	//	}
	//	// start scaping the page under the link found
	//	e.Request.Visit(link)
	//})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// On every <a> element with collection-product-card class call callback
	//c.OnHTML(`a.collection-product-card`, func(e *colly.HTMLElement) {
	//	// Activate detailCollector if the link contains "coursera.org/learn"
	//	courseURL := e.Request.AbsoluteURL(e.Attr("href"))
	//	if strings.Index(courseURL, "coursera.org/learn") != -1 {
	//		detailCollector.Visit(courseURL)
	//	}
	//})

	// Extract details of the course
	//detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
	//	log.Println("Course found", e.Request.URL)
	//	title := e.ChildText(".banner-title")
	//	if title == "" {
	//		log.Println("No title found", e.Request.URL)
	//	}
	//	course := Course{
	//		Title:       title,
	//		URL:         e.Request.URL.String(),
	//		Description: e.ChildText("div.content"),
	//		Creator:     e.ChildText("li.banner-instructor-info > a > div > div > span"),
	//		Rating:      e.ChildText("span.number-rating"),
	//	}
	//	// Iterate over div components and add details to course
	//	e.ForEach(".AboutCourse .ProductGlance > div", func(_ int, el *colly.HTMLElement) {
	//		svgTitle := strings.Split(el.ChildText("div:nth-child(1) svg title"), " ")
	//		lastWord := svgTitle[len(svgTitle)-1]
	//		switch lastWord {
	//		// svg Title: Available Languages
	//		case "languages":
	//			course.Language = el.ChildText("div:nth-child(2) > div:nth-child(1)")
	//		// svg Title: Mixed/Beginner/Intermediate/Advanced Level
	//		case "Level":
	//			course.Level = el.ChildText("div:nth-child(2) > div:nth-child(1)")
	//		// svg Title: Hours to complete
	//		case "complete":
	//			course.Commitment = el.ChildText("div:nth-child(2) > div:nth-child(1)")
	//		}
	//	})
	//	courses = append(courses, course)
	//})

	// Start scraping on http://coursera.com/browse
	c.Visit("https://18f.github.io/API-All-the-X/")

	//enc := json.NewEncoder(file)
	//enc.SetIndent("", "  ")

	// Dump json to the standard output
	//enc.Encode(courses)
}
