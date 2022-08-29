package reader

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	linkparser "github.com/sadex11/gopher04parser"
)

func getPageData(pageUrl string) ([]byte, error) {
	resp, err := http.Get(pageUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func processPage(pageUrl string, linkOut chan<- string) {
	log.Println("Processing: ", pageUrl)
	content, err := getPageData(pageUrl)

	if err != nil {
		log.Println("Process url error", err)
		return
	}

	for _, link := range *linkparser.GetNodeLinks(bytes.NewReader(content)) {
		linkOut <- link.Href
	}
}

func getFullUrl(baseUrl string, refUrl string) string {
	// TODO
	return baseUrl + "/" + refUrl
}

func GetSiteStructure(baseUrl string, level int) {
	// TODO use different approach to channel -> wg
	log.SetOutput(os.Stdout)

	baseContent, err := getPageData(baseUrl)

	if err != nil {
		panic(err)
	}

	if baseContent == nil {
		log.Println("Empty base content")
		return
	}

	var processLinks []string

	for _, link := range *linkparser.GetNodeLinks(bytes.NewReader(baseContent)) {
		processLinks = append(processLinks, link.Href)
	}

	if len(processLinks) == 0 {
		log.Println("No base links to process")
		return
	}

	linkChannel := make(chan string)
	// TODO check links if contains base
	for _, refUrl := range processLinks {
		go processPage(getFullUrl(baseUrl, refUrl), linkChannel)
	}

	remains := len(processLinks)

	for remains > 0 {
		select {
		case newLink := <-linkChannel:
			log.Println("New link", newLink)
			remains--
			log.Println("Remains", remains)
			remains++
		}
	}
}
