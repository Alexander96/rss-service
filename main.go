package main

import (
	"fmt"

	rssreader "github.com/Alexander96/rssreader"
)

func main() {
	urls := []string{
		"http://rss.cnn.com/rss/edition.rss",
		"http://rss.cnn.com/rss/edition_world.rss",
		"http://rss.cnn.com/rss/edition_travel.rss",
		"http://rss.cnn.com/rss/cnn_latest.rss"}

	items := rssreader.Parse(urls)

	fmt.Println(items)
}
