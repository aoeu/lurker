package lurker

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	maestroURL = "http://192.168.34.10:3000"
	userAgent  = "ComicGator (https://github.com/comicgator/lurker)"
	lurkerUser = "c11z"
	lurkerPass = "c11z"
)

var (
	client *http.Client
)

func init() {
	client = &http.Client{}
}

func FetchComicPage(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func SaveResource(id, resource string, body []byte) error {
	var res *http.Response
	var err error

	// Try post first
	var postURL string = maestroURL + "/" + resource
	res, err = maestro(postURL, "POST", body)
	if err != nil {
		return err
	}
	// If successful return
	if res.StatusCode == 201 {
		return nil
	}
	// Otherwise Patch by filtering against the unique id
	var patchURL string = postURL + "?id=eq." + id
	res, err = maestro(patchURL, "PATCH", body)
	if err != nil {
		return err
	}
	// If successful return
	if res.StatusCode == 204 {
		return nil
	}
	// Oops something has gone horribly wrong
	log.Println("oops something went horribly wrong " + res.Status)
	return errors.New("Unable to post or patch " + resource + " with id= " + id)
}

func maestro(url, verb string, body []byte) (res *http.Response, err error) {
	var req *http.Request
	req, err = http.NewRequest(verb, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", userAgent)
	req.SetBasicAuth(lurkerUser, lurkerPass)
	res, err = client.Do(req)
	return
}

// TODO: Change name to be more specific
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
