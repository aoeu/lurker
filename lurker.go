// lurk file contains control flow for different crawling strategies.
package lurker

import (
	"log"
)

// Full function gets all strips from all comics starting at the first page.
func ETL(comicList []Comic, delta bool) {
	log.Println("Running ETL.")
	// Upsert comics to database
	if delta { 
		
	}
	// for each Comic spawn goroutine, crawl and return []Strip
	// TODO: goroutine
	for _, comic := range comicList {
		var count int
		var startUrl string
		if delta {
			log.Println("Delta requested.")
			// Go get last Strip for each comic from the database
			// Set Comic.FirstPageUrl to last Strip.Url
			// Set Comic.StripCount to Strip.Number
			// TODO: remove temp setting of count and startUrl
			count = 1
			startUrl = comic.FirstPageUrl
		} else {
			count = 1
			startUrl = comic.FirstPageUrl
		}

		var strips []Strip
		strips = Crawl(comic, startUrl, count)
		for _, strip := range strips {
			log.Print(string(strip.Export()))
		}
		// Validate and Upsert Strips to database
	}
}
