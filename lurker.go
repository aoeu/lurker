// lurk file contains control flow for different crawling strategies.
package lurker

import (
	"log"
	"net/url"

)

// Full function gets all strips from all comics starting at the first page.
func ETL(comicList []Comic, delta bool) {
	log.Println("Running ETL")
	// Upsert comics to database
	// for each Comic spawn goroutine, crawl and return []Strip
	// TODO: goroutine
	for _, comic := range comicList {
		var count int
		var start *url.URL
		if delta {
			log.Println("Delta requested")
			// Go get last Strip for each comic from the database
			// Set Comic.FirstPageUrl to last Strip.Url
			// Set Comic.StripCount to Strip.Number
			// TODO: remove temp setting of count and startUrl
			count = 0
			start, _ = url.Parse(comic.FirstPageUrl)
			// if err != nil {
				// panic(err)
			// }
			// start = last
		} else {
			count = 0
			start, _ = url.Parse(comic.FirstPageUrl)
			
		}
		log.Printf("Crawling Comic %s\n", comic.Hostname)
		log.Printf("Starting at url %s\n", start.String())
		strips := Crawl(comic, start, count)
		log.Println("Crawling Complete")
		for _, strip := range strips {
			log.Print(string(strip.Export()))
		}
		// Validate and Upsert Strips to database
	}
}
