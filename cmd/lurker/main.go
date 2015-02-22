package main

import (
	"fmt"
	_ "github.com/moovweb/gokogiri"
	_ "github.com/moovweb/gokogiri/html"
	_ "github.com/moovweb/gokogiri/xpath"

	"github.com/comicgator/lurker"
	// "github.com/comicgator/lurker/comic"
)

func main() {
	fmt.Println("Hello World")

	comics := lurker.LoadComics()
	// for each comic make put request to maestro
	for _, comic := range comics {
		// id := maestro.Put(comic.Export())	
		// comic.SetId(id)
		fmt.Printf("%+v\n", comic)
	}
}

