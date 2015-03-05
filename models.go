// Comic file contains interface for comic structures and defines methods
// for parsing comicstrip webpage html.
package lurker

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	_ "log"
	_ "reflect"
)

var (
	namespace = uuid.NamespaceOID
)

type Comic struct {
	ID               string   `json:"-"`
	Hostname         string   `json:"hostname"`
	Title            string   `json:"title"`
	Creator          string   `json:"creator"`
	HeadlineImageURL string   `json:"headline_image_url"`
	FirstPageURL     string   `json:"first_page_url"`
	Pattern          *Pattern `json:"pattern"`
}

type Pattern struct {
	Method     string   `json:"method"`
	Title      []string `json:"title"`
	Prev       []string `json:"prev"`
	Next       []string `json:"next"`
	Image      []string `json:"image"`
	AltText    []string `json:"alt_text"`
	BonusImage []string `json:"bonus_image"`
}

func (c *Comic) GenerateUUID() {
	c.ID = uuid.NewV5(namespace, c.Hostname).String()
}

func (c Comic) Save() {
	SaveResource(c.ID, "comic", c.Export())
}

func (c Comic) Export() []byte {
	m := map[string]interface{}{} // ideally use make with the right capacity
	m["id"] = c.ID
	m["hostname"] = c.Hostname
	m["title"] = c.Title
	m["creator"] = c.Creator
	m["headline_image_url"] = c.HeadlineImageURL
	output, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return output
}

type Strip struct {
	ID                string `json:"id"`
	ComicId           string `json:"comic_id"`
	Title             string `json:"title"`
	Number            int    `json:"number"`
	URL               string `json:"url"`
	ImageURL          string `json:"image_url"`
	ThumbnailImageURL string `json:"thumbnail_image_url"`
	BonusImageURL     string `json:"bonus_image_url"`
	AltText           string `json:"alt_text"`
}

func (s *Strip) GenerateUUID() {
	s.ID = uuid.NewV5(namespace, s.URL).String()
}

func (s Strip) Save() {
	SaveResource(s.ID, "strip", s.Export())
}

func (s Strip) Export() []byte {
	// basic template for exporting Strip as json
	output, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return output
}
