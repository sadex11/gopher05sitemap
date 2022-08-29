package main

import (
	"github.com/sadex11/gopher05sitemap/reader"
)

func main() {
	reader.GetSiteStructure("http://127.0.0.1:8080", 0)
}
