// Comic file contains interface for comic structures and defines methods
// for parsing comicstrip webpage html.
package lurker

import (
	"encoding/json"
	_ "log"
	_ "reflect"
)

type Comic struct {
	id               string   `json:"-"`
	Url              string   `json:"url"`
	Hostname         string   `json:"hostname"`
	Title            string   `json:"title"`
	Creator          string   `json:"creator"`
	HeadlineImageUrl string   `json:"headline_image_url"`
	FirstPageUrl     string   `json:"first_page_url"`
	Pattern          *Pattern `json:"pattern"`
}

type Pattern struct {
	Method     string   `json:"method"`
	Title      []string `json:"title"`
	Prev       []string `json:"prev"`
	Next       []string `json:"next"`
	Image      []string `json:"image"`
	BonusImage []string `json:"bonus_image"`
}

func (c *Comic) Export() []byte {
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

type Strip struct {
	Id                string `json:"-"`
	ComicId           string `json:"comic_id"`
	Title             string `json:"title"`
	Number            int    `json:"number"`
	Url               string `json:"url"`
	ImageUrl          string `json:"image_url"`
	ThumbnailImageUrl string `json:"thumbnail_image_url"`
	BonusImageUrl     string `json:"bonus_image_url"`
	AltText           string `json:"alt_text"`
}

func (s Strip) Export() []byte {
	// basic template for exporting Strip as json
	output, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return output
}
