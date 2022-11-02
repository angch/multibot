package sdapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Server struct {
	URL *url.URL
}

func NewServer(host string) *Server {
	u, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &Server{URL: u}
}

func (s *Server) Txt2Img(prompt string) ([]byte, error) {
	// quick hack
	p := NewTxt2ImgParameters()
	p.Prompt = prompt

	u := s.URL.String()
	u += "/sdapi/v1/txt2img"
	log.Println(u)
	resp, err := http.Post(u, "application/json", p.IoReader())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// FIXME: check http errorcode, etc

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := Txt2ImgParametersResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(result.Images) < 1 {
		return nil, err
	}
	first := result.Images[0]
	// hack, hack
	// "data:image/png;base64,xxxx"
	if strings.HasPrefix(first, "data:image/png;base64,") {
		first = first[len("data:image/png;base64,"):]
		image, err := base64.StdEncoding.DecodeString(first)
		if err != nil {
			return nil, err
		}
		return image, nil
	}

	return nil, fmt.Errorf("No image")
}
