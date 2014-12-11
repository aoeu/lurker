package comic

import (
        "encoding/json"
        // "fmt"
)

type Xpath struct {
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

func (x Xpath) Prev() string {
	//temp return pattern
	return x.pattern.prev
}

func (x Xpath) Next() string {
	//temp return pattern
	return x.pattern.next
}

func (x Xpath) Image() string {
	//temp return pattern
	return x.pattern.image
}

func (x Xpath) AltText() string {
	//temp return pattern
	return x.pattern.altText
}

func (x Xpath) BonusImage() string {
	//temp return pattern
	return x.pattern.bonusImage
}
