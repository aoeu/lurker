// lurk file contains control flow for different crawling strategies.
package lurker

import (
	"log"
)

// Full function gets all strips from all comics starting at the first page.
func ETL(comicList []Comic, delta bool) {
	log.Println("Running ETL")
	// Upsert comics to database
	// for each Comic spawn goroutine, crawl and return []Strip
	// TODO: goroutine
	for _, comic := range comicList {
		comic.GenerateUUID()
		comic.Save()
		var count int
		var start string
		if delta {
			log.Println("Delta requested")
			// Go get last Strip for each comic from the database
			// Set Comic.FirstPageURL to last Strip.URL
			// Set Comic.StripCount to Strip.Number
			// TODO: remove temp setting of count and startURL
			count = 0
			start = comic.FirstPageURL
			// if err != nil {
			// panic(err)
			// }
			// start = last
		} else {
			count = 0
			start = comic.FirstPageURL

		}
		log.Printf("Crawling Comic %s\n\n", comic.Hostname)
		log.Print(string(comic.Export()))
		strips := Crawl(comic.ID, comic.Hostname, start, count, comic.Pattern)
		log.Println("Crawl Complete")
		for _, strip := range strips {
			log.Print(string(strip.Export()))
			strip.Save()
		}
		// Validate and Upsert Strips to database
	}
}
