package main

import (
	"fmt"
	_ "github.com/PuerkitoBio/fetchbot"
	_ "github.com/PuerkitoBio/purell"
	_ "github.com/moovweb/gokogiri"
	_ "github.com/moovweb/gokogiri/html"
	_ "github.com/moovweb/gokogiri/xpath"

	_ "github.com/c11z/lurker"
	// _ "github.com/c11z/lurker/thumb"
	"github.com/c11z/lurker/comic"
)

func main() {
	fmt.Println("Hello World")

	comics := comic.Load()
	// for each comic make put request to maestro
	for _, comic := range comics {
		// id := maestro.Put(comic.Export())	
		// comic.SetId(id)
		fmt.Printf("%+v\n", comic)
	}
}

