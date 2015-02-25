// Comic file contains interface for comic structures and defines methods
// for parsing comicstrip webpage html.
package lurker

import (
	"log"
	"encoding/json"
	"github.com/moovweb/gokogiri"
	_ "github.com/moovweb/gokogiri/html"
	_ "github.com/moovweb/gokogiri/xpath"
	_ "reflect"
)

type Comic struct {
	id               string `json:"-"`
	Hostname         string `json:"hostname"`
	Title            string `json:"title"`
	Creator          string `json:"creator"`
	HeadlineImageUrl string `json:"headline_image_url"`
	FirstPageUrl     string `json:"first_page_url"`
	Pattern          struct {
		Method     string `json:"method"`
		Title      string `json:"title"`
		Prev       string `json:"prev"`
		Next       string `json:"next"`
		Image      string `json:"image"`
		BonusImage string `json:"bonus_image"`
	} `json:"pattern"`
}

func (c Comic) Export() []byte {
	m := map[string]interface{}{} // ideally use make with the right capacity
	m["hostname"] = c.Hostname
	m["title"] = c.Title
	m["creator"] = c.Creator
	m["headline_image_url"] = c.HeadlineImageUrl
	output, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return output
}

func (c *Comic) SetId(id string) {
	c.id = id
}

func (c Comic) Id() string {
	return c.id
}

func (c *Comic) Parse(page []byte) (Strip, string)  {
	var strip Strip = Strip{}
	var nextUrl string = ""

	strip.Comic = c
	strip.ComicId = c.Id()


	doc, err := gokogiri.ParseHtml(page)
	if err != nil {
		log.Fatalln(err)
	}
	defer doc.Free() 
	
	// Get title from doc
	// xtitle := xpath.Compile(c.Pattern.Title)
	// stitle, err := doc.Root().Search(xtitle)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// title := stitle[0].Content()
	// strip.Title = title

	// Don't need previous at the moment
	// xprev := xpath.Compile(c.Pattern.Prev)
	// sprev, err := doc.Root().Search(xprev)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// prev := sprev[0].Attr("href")

	// xnext := xpath.Compile(c.Pattern.Next)
	// snext, err := doc.Root().Search(xnext)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// next := snext[0].Attr("href")
	// if next != "" {
	// 	nextUrl = "http://" + c.Hostname + next	
	// }
	

	// ximage := xpath.Compile(c.Pattern.Image)
	// simage, err := doc.Root().Search(ximage)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// imageUrl := simage[0].Attr("src")
	// strip.ImageUrl = imageUrl

	// altText := simage[0].Attr("title")
	// strip.AltText = altText

	// bonus image is causing memory leak
	// xbonusimage := xpath.Compile(c.Pattern.BonusImage)
	// sbonusimage, _ := doc.Root().Search(xbonusimage)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// bonusImage := sbonusimage[0].Attr("href")

	// log.Println(bonusImage)
	// log.Println(reflect.TypeOf(xbonusimage))

	return strip, nextUrl
}



type Strip struct {
	Comic             *Comic
	Id                string `json:"-"`
	ComicId           string `json:comic_id`
	Title             string `json:title`
	Number            int    `json:number`
	Url               string `json:url`
	ImageUrl          string `json:image_url`
	ThumbnailImageUrl string `json:thumbnail_image_url`
	BonusImageUrl     string `json:bonus_image_url`
	AltText           string `json:alt_text`
}

func (s Strip) Export() []byte {
	// basic template for exporting Strip as json
	output, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return output
}
