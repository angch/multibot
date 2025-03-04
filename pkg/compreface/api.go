package compreface

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func getapi(u string, apikey string) ([]byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("x-api-key", apikey)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Body == nil {
		log.Println(err)
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp.Body.Close()
	return data, nil
}

func postapi(u string, apikey string, postdata []byte) ([]byte, error) {
	postbody := bytes.NewBuffer(postdata)
	req, err := http.NewRequest("POST", u, postbody)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("x-api-key", apikey)
	req.Header.Add("Content-Type", "multipart/form-data")
	log.Println("Posting to ", u, apikey, len(postdata))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.Body == nil {
		log.Println(err)
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp.Body.Close()
	return data, nil
}
