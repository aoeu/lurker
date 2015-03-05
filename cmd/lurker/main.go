package main

import (
	"encoding/json"
	"github.com/comicgator/lurker"
	"io/ioutil"
	"log"
)

func main() {
	log.Println("Lurker Up")
	var delta bool = false

	// Read in comics.json file
	log.Println("Loading comics...")
	comics := LoadComics()
	log.Printf("%d Comics Loaded\n", len(comics))
	// TODO: insert command line option logic to call specific running
	// strategies.
	lurker.ETL(comics, delta)
}

// Reads json file and unmarshals into list of comic structs.
func LoadComics() []lurker.Comic {
	var filename string = "/home/vagrant/go/src/github.com/comicgator/lurker/comics.json"
	// var models []model
	var comics []lurker.Comic
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	// Unmarshal into unopinonated model.
	err = json.Unmarshal(source, &comics)
	if err != nil {
		panic(err)
	}
	return comics
}
