package lurker

import (
	"errors"
	"github.com/PuerkitoBio/purell"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"log"
	"strings"
	"time"
)

func Crawl(comicID, hostname, url string, count int, pattern *Pattern) []Strip {
	var strips []Strip
	visited := make(map[string]bool)
	for visited[url] != true {
		count += 1
		var newStrip Strip

		newStrip.URL = url
		newStrip.GenerateUUID()
		newStrip.ComicId = comicID
		newStrip.Number = count

		log.Println("Fetching url " + url)
		page := FetchComicPage(url)
		log.Println("Parsing...")
		// Every comic should have a strip image url.
		image, err := parseImage(page, pattern.Image)
		if err != nil {
			log.Println(err)
		} else {
			newStrip.ImageURL = buildURL(hostname, image)
		}

		if len(pattern.Title) > 0 {
			title, err := parseTitle(page, pattern.Title)
			if err != nil {
				log.Println(err)
			} else {
				newStrip.Title = title
			}
		}

		if len(pattern.AltText) > 0 {
			altText, err := parseAltText(page, pattern.Image)
			if err != nil {
				log.Println(err)
			} else {
				newStrip.AltText = altText
			}
		}

		if len(pattern.BonusImage) > 0 {
			bonus, err := parseImage(page, pattern.BonusImage)
			if err != nil {
				log.Println(err)
			} else {
				newStrip.BonusImageURL = buildURL(hostname, bonus)
			}
		}

		strips = append(strips, newStrip)

		// Get next url to crawl. If there is an error finding the next
		// endpoint then break the loop.
		endpoint, err := parseNext(page, pattern.Next)
		if err != nil {
			log.Println(err)
			break
		} else {
			// Store url in visited urls and then reassign to new url
			visited[url] = true
			newURL := buildURL(hostname, endpoint)
			url = newURL
		}

		// Sleep 5 seconds to limit abuse to comic websites
		log.Println("Parse Complete!")
		log.Println("Resting...\n")
		time.Sleep(5 * time.Second)
	}
	return strips
}

func buildURL(hostname, input string) string {
	// Assumes that next link is always and endpoint
	// TODO: be smarter about building the url
	var raw string
	// Key it simple and check if url is relative or not
	if strings.Contains(input, "http") {
		raw = input
	} else {
		raw = "http://" + hostname + input
	}
	url, err := cleanURL(raw)
	if err != nil {
		log.Println(raw)
		log.Fatalln(err)
	}
	return url
}

func cleanURL(input string) (output string, err error) {
	output, err = purell.NormalizeURLString(input, purell.FlagsUnsafeGreedy)
	return
}

func searchXpath(page []byte, paths []string) []xml.Node {
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
	nodes = searchXpath(page, paths)
	for _, node := range nodes {
		attrCheck := node.Attribute("href")
		if attrCheck != nil {
			value := node.Attr("href")
			if value != "" && value != "#" {
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

func parseImage(page []byte, paths []string) (url string, err error) {
	var sources []string

	nodes := searchXpath(page, paths)
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
		url, err = cleanURL(item)
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

	nodes := searchXpath(page, paths)
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
	if len(alts) > 0 {
		alt = alts[0]
	} else {
		err = errors.New("Unable to find alt text.")
	}
	return
}

func parseTitle(page []byte, paths []string) (title string, err error) {
	var titles []string
	nodes := searchXpath(page, paths)
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
