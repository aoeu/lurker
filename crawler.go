package lurker

import (
	"errors"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	// "reflect"
)

func Crawl(comic Comic, uri *url.URL, count int) []Strip {
	var strips []Strip
	visited := make(map[*url.URL]bool)
	for visited[uri] != true {
		count += 1
		var newStrip Strip
		newStrip.ComicId = comic.Id()
		newStrip.Url = uri.String()
		newStrip.Number = count

		page := fetch(uri)
		title, err := parseTitle(page, comic.Pattern.Title)
		if err != nil {
			log.Println(err)
		} else {
			newStrip.Title = title
		}
		image, err := parseImage(page, comic.Pattern.Image)
		if err != nil {
			log.Println(err)
		} else {
			newStrip.ImageUrl = image.String()
		}
		// Alt text is retrieved from the same node as the Image thus
		// uses the same pattern
		altText, err := parseAltText(page, comic.Pattern.Image)
		if err != nil {
			log.Println(err)
		} else {
			newStrip.AltText = altText
		}
		bonus, err := parseImage(page, comic.Pattern.BonusImage)
		if err != nil {
			log.Println(err)
		} else {
			newStrip.BonusImageUrl = bonus.String()
		}

		strips = append(strips, newStrip)

		// Get next uri to crawl
		endpoint, err := parseNext(page, comic.Pattern.Next)
		if err != nil {
			// There is an error finding the next endpoint break
			// the loop
			log.Println(err)
			break
		} else {
			// Store uri in visited urls and then reassign to new uri
			visited[uri] = true
			newUri, err := url.Parse(comic.Url + endpoint)
			if err != nil {
				// The new endpoint doesn't make sense, break
				// the loop
				break
			}
			// everyrthing is good assign newUri to uri and loop again
			uri = newUri
		}

		// Sleep 5 seconds to limit abuse to comic websites
		log.Println("Sleepy time\n\n")
		time.Sleep(5 * time.Second)
	}
	return strips
}

func fetch(uri *url.URL) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri.String(), nil)
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

func xpathParse(page []byte, paths []string) []xml.Node {
	// Parse html page
	doc, err := gokogiri.ParseHtml(page)
	if err != nil {
		log.Fatalln(err)
	}
	defer doc.Free()

	var results []xml.Node
	// Some comics might have different html structure between pages, here
	// we try several paths and collect the resultant Nodes found.
	for _, path := range paths {
		expression := xpath.Compile(path)
		nodeList, err := doc.Root().Search(expression)
		if err != nil {
			log.Fatalln(err)
		}
		// log.Println(reflect.TypeOf(nodeList))
		for _, node := range nodeList {
			results = append(results, node.Duplicate(1))
		}
		expression.Free()
	}
	return results
}

func parseNext(page []byte, paths []string) (next string, err error) {
	var endpoints []string
	var nodes []xml.Node
	nodes = xpathParse(page, paths)
	for _, node := range nodes {
		
		attrCheck := node.Attribute("href")
		if attrCheck != nil {
			value := node.Attr("href")
			log.Println("Endpoint value: " + value)
			if value != "" || value != "#" {
				endpoints = append(endpoints, value)
			}
		}
	}
	if len(endpoints) > 0 {
		next = endpoints[0]
	} else {
		err = errors.New("Unable to find next endpoint.")
	}
	return
}

func parseImage(page []byte, paths []string) (uri *url.URL, err error) {
	var sources []string

	nodes := xpathParse(page, paths)
	for _, node := range nodes {
		srcCheck := node.Attribute("src")
		if srcCheck != nil {
			value := node.Attr("src")
			if value != "" {
				sources = append(sources, value)
			}
		}
	}
	if len(sources) > 0 {
		item := sources[0]
		uri, err = url.Parse(item)
		if err != nil {
			err = errors.New("Unable to parse url " + item)
		}

	} else {
		err = errors.New("Unable to find image source.")
	}
	return
}

func parseAltText(page []byte, paths []string) (alt string, err error) {
	var alts []string

	nodes := xpathParse(page, paths)
	for _, node := range nodes {
		titleCheck := node.Attribute("title")
		altCheck := node.Attribute("alt")
		if titleCheck != nil {
			value := node.Attr("title")
			if value != "" {
				alts = append(alts, value)
			}

		} else if altCheck != nil {
			value := node.Attr("alt")
			if value != "" {
				alts = append(alts, value)
			}
		}
	}
	if len(alts) < 0 {
		alt = alts[0]
	} else {
		err = errors.New("Unable to find alt text.")
	}
	return
}

func parseTitle(page []byte, paths []string) (title string, err error) {
	var titles []string
	nodes := xpathParse(page, paths)
	for _, node := range nodes {
		value := node.Content()
		if value != "" {
			titles = append(titles, value)
		}
	}
	if len(titles) > 0 {
		title = titles[0]
	} else {
		err = errors.New("Unable to find page title.")
	}
	return
}
