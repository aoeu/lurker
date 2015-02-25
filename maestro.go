// maestro package handles all requests to api machine
package lurker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Put(reqBody []byte) string {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", "http://192.168.34.10:3000/comics", bytes.NewBuffer(reqBody))
	req.Header.Set("Accepts", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	log.Println("response Status:", res.Status)
	log.Println("response Headers:", res.Header)
	resBody, _ := ioutil.ReadAll(res.Body)
	log.Println("response Body:", string(resBody))
	var comic map[string]string

	err = json.Unmarshal(resBody, &comic)
	if err != nil {
		panic(err)
	}

	return comic["comic_id"]
}