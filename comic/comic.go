// Comic package contains interface for comic structures and defines methods
// for parsing panel webpage html.
package comic

import (
	// "fmt"
	"io/ioutil"
	"encoding/json"
)

type model map[string]interface{}

type Comic interface {
	Export() []byte
	Parse() Panel
	Id() string
	SetId(string)
}

type Panel struct {
	Comic *Comic
	Id int
	PreviousPanel *Panel
	NextPanel *Panel
	Title string
	Number int
	Url string
	ImageUrl string
	ThumbnailImageUrl string
	BonusImageUrl string
	AltText string
}

func (p Panel) Export() []byte {
	// basic template for exporting Panel as json
	// TODO: actually fill data from Panel struct 
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

func Load() []Comic {
	var filename string = "../../comics.json"
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
					prev: p["prev"].(string),
					next: p["next"].(string),
					image: p["image"].(string),
					altText: p["alt_text"].(string),
					bonusImage: p["bonus_image"].(string)}}
			comics = append(comics, xpath)
		}
		// fmt.Printf("%+v\n", m)
	}
	return comics
}