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
	Prev() string
	Next() string
	Image() string
	AltText() string
	BonusImage() string
}

func Load() []Comic {
	var filename string = "comics.json"
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
				pattern: Pattern{
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