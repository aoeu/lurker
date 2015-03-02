package lurker

import (
	"log"
	"net/http"
	"io/ioutil"
)
const (
	maestroURL = "http://192.168.34.10:3000"
	userAgent = "ComicGator (https://github.com/comicgator/lurker)"
)

var (
	client *http.Client
)

func init() {
	client = &http.Client{}
}

func FetchComicPage(url string) []byte {
	body, err := get(url)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}



func SaveComic(hostname string, export []byte) {
}

func post(url string, body []byte) (res *http.Response, err error) {
	return
}

func patch(url string, body []byte) (res *http.Response, err error) {
	return
}


func get(url string) (body []byte, err error) {
	var req *http.Request
	var res *http.Response

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", userAgent)

	res, err = client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}

// func Save(payload []byte, endpoint string) (string, error) {
// 	// Try POST

// 	// if http 400 then STATUS 400 Bad Request
// TIME 38 ms
// Pretty Raw Preview  JSON Copy
// {
// "hint": null,
// "details": "Key (url)=(http://xkcd.com/1488) already exists.",
// "code": "23505",
// "message": "duplicate key value violates unique constraint \"strip_uq_url\""
// }

// Location → /comic?id=eq.35a36d33-3c62-4336-a075-09ad96942432




// }
