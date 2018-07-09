package requests

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const OK = "200 OK"
const HTTP = "http://"

func MakeJsonHTTPPOSTRequest(url string, json []byte) (string, []byte, error) {
	req, err := http.NewRequest("POST", "http://"+url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", []byte{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.Status, body, nil
}

func MakeHTTPGetRequest(url string) (string, []byte, error) {
	response, err := http.Get("http://" + url)
	if err != nil {
		return "", []byte{}, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	return response.Status, content, err
}
