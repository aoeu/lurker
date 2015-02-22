// Comic file contains interface for comic structures and defines methods
// for parsing strip webpage html.
package lurker

import (
	// "fmt"
	"io/ioutil"
	"encoding/json"
)

type model map[string]interface{}

type Comic interface {
	Export() []byte
	Parse() Strip
	Id() string
	SetId(string)
}

type Strip struct {
	Comic *Comic
	Id int
	PreviousStrip *Strip
	NextStrip *Strip
	Title string
	Number int
	Url string
	ImageUrl string
	ThumbnailImageUrl string
	BonusImageUrl string
	AltText string
}

func (p Strip) Export() []byte {
	// basic template for exporting Strip as json
	// TODO: actually fill data from Strip struct 
	panel := struct {
		ComicId int `json:comic_id`
		PreviousPanelId int `json:previous_panel_id`
		NextPanelId int `json:next_panel_id`
		Title string `json:title`
		Number int `json:number`
		Url string `json:url`
		ImageUrl string `json:image_url`
		ThumbnailImageUrl string `json:thumbnail_image_url`
		BonusImageUrl string `json:bonus_image_url`
		AltText string `json:alt_text`
	} {
		0,
		0,
		0,
		"first panel",
		1,
		"http://example.com/1",
		"http://example.com/1.png",
		"http://comicgator.com/panel.png",
		"",
		"First panel",
	}

	output, err := json.Marshal(panel)
	if err != nil {
		panic(err)
	}
	return output
}


type Xpath struct {
	id		 string    `json:"-"`
	Hostname         string `json:"hostname"`
	Title            string `json:"title"`
	Creator          string `json:"creator"`
	HeadlineImageUrl string `json:"headline_image_url"`
	pattern          Pattern `json:"-"`
}

type Pattern struct {
	title      string
	prev       string
	next       string
	image      string
	bonusImage string
}

func (x Xpath) Export() []byte {
	output, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return output
}

func (x Xpath) Parse() Strip {
	//temp return pattern
	return Strip {
		nil,
		0,
		nil,
		nil,
		"first panel",
		1,
		"http://example.com/1",
		"http://example.com/1.png",
		"http://comicgator.com/panel.png",
		"",
		"First panel",
	}
}

func (x Xpath) Id() string {
	return x.id
}

func (x Xpath) SetId(id string) {
	x.id = id
}

func LoadComics() []Comic {
	var filename string = "/home/vagrant/go/src/github.com/comicgator/lurker/comics.json"
	var models []model
	var comics []Comic
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(source, &models)
	if err != nil {
		panic(err)
	}

	for _, m := range models {
		if m["method"] == "xpath" {
			p := m["pattern"].(map[string]interface{})
			xpath := Xpath {
				Hostname: m["hostname"].(string),
				Title: m["title"].(string),
				Creator: m["creator"].(string),
				HeadlineImageUrl: m["headline_image_url"].(string),
				pattern: Pattern {
					title: p["title"].(string),
					prev: p["prev"].(string),
					next: p["next"].(string),
					image: p["image"].(string),
					bonusImage: p["bonus_image"].(string)}}
			comics = append(comics, xpath)
		}
		// fmt.Printf("%+v\n", m)
	}
	return comics
}