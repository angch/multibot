package sdapi

import (
	"encoding/base64"
	"encoding/json"
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

var negativePromptsArray = []string{
	"(disfigured)",
	"(deformed)",
	"(poorly drawn)",
	"(extra limbs)",
	"boring",
	"sketch",
	"lackluster",
	"signature",
	"letters",
	"watermark",
	"low res",
	"horrific",
	"mutated",
	"artifacts",
	"bad art",
	"gross",
	"poor quality",
	"low quality",
}

var negativePrompt string

var enhancedPrompts = []string{
	// "dark and gloomy",
	// "full body",
	"8k unity render",
	"skin pores",
	"detailed iris",
	// "very dark lighting",
	// "heavy shadows",
	"detailed",
	"detailed face",
	"(vibrant)",
	// "photo realistic",
	// "realistic",
	// "dramatic",
	// "dark",
	"sharp focus",
	"(8k)",
}

var enhancedPrompt string

func init() {
	negativePrompt = strings.Join(negativePromptsArray, ", ")
	enhancedPrompt = strings.Join(enhancedPrompts, ", ")
}

func (s *Server) Txt2Img(prompt string) ([]byte, error) {
	// quick hack
	p := NewTxt2ImgParameters()
	p.Prompt = prompt

	if len(prompt) < 20 {
		p.Prompt = p.Prompt + ", " + enhancedPrompt
	}
	p.NegativePrompt = negativePrompt
	p.SetSampler("DDIM")

	u := s.URL.String()
	u += "/sdapi/v1/txt2img"
	log.Println(u, p.IoReader().String())
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
		log.Println(string(body))
		return nil, err
	}
	first := result.Images[0]
	// hack, hack
	// "data:image/png;base64,xxxx"
	first = strings.TrimPrefix(first, "data:image/png;base64,")

	image, err := base64.StdEncoding.DecodeString(first)
	if err != nil {
		return nil, err
	}
	return image, nil
	// log.Println(string(body))
	// log.Printf("%+v\n", result)

	// return nil, fmt.Errorf("No image")
}
