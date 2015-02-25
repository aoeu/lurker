package lurker

import (
	"log"
	"net/http"
	"time"
	"io/ioutil"
)

func Crawl(comic Comic, url string, count int) []Strip {
	var strips []Strip
	for url != "" {
		page := fetch(url)
		newStrip, newUrl := comic.Parse(page)
		url = newUrl
		newStrip.Number = count
		strips = append(strips, newStrip)
		count += 1
		// Sleep 5 seconds to limit abuse to comic websites
		log.Println("Sleepy time")
		time.Sleep(5 * time.Second)
	}
	return strips
}

func fetch(url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	
	req.Header.Set("User-Agent", "ComicGatorBot/1.0 (https://github.com/comicgator/lurker)")

        resp, err := client.Do(req)
        if err != nil {
                log.Fatalln(err)
        }

        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatalln(err)
        }
        return body
}
