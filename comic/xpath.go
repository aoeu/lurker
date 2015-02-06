package comic

import (
	"encoding/json"
	// "fmt"
)

type Xpath struct {
	id		 string    `json:"-"`
	Hostname         string `json:"hostname"`
	Title            string `json:"title"`
	Creator          string `json:"creator"`
	HeadlineImageUrl string `json:"headline_image_url"`
	pattern          Pattern `json:"-"`
}

type Pattern struct {
	prev       string
	next       string
	image      string
	altText    string
	bonusImage string
}

func (x Xpath) Export() []byte {
	output, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return output
}

func (x Xpath) Parse() Panel {
	//temp return pattern
	return Panel {
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