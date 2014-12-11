package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"lurker/comic"
	_ "lurker/thumb"
	"net/http"
)

func callMaestro(reqBody []byte) {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", "http://192.168.34.10:3000/comics", bytes.NewBuffer(reqBody))
	req.Header.Set("Accepts", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)
	fmt.Println("response Headers:", res.Header)
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println("response Body:", string(resBody))
}

func main() {
	fmt.Println("Hello World")

	var comics []comic.Comic
	comics = comic.Load()

	for _, comic := range comics {
		callMaestro(comic.Export())
	}
}
