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
	"time"
)

type Server struct {
	URL *url.URL

	NegativePrompt string
	EnhancedPrompt string
	Models         map[string]Model
	Config         *Config
}

func NewServer(host string) *Server {
	u, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return nil
	}
	s := Server{URL: u}
	s.NegativePrompt = strings.Join(negativePromptsArray, ", ")
	s.EnhancedPrompt = strings.Join(enhancedPrompts, ", ")
	s.Models = make(map[string]Model, 0)

	go func() {
		s.Config, _ = s.GetConfig()
	}()
	return &s
}

var negativePromptsArray = []string{
	"(disfigured)",
	"(deformed)",
	"(poorly drawn)",
	"(extra limbs)",
	// "boring",
	"sketch",
	// "lackluster",
	"signature",
	"letters",
	"watermark",
	"low res",
	"horrific",
	"mutated",
	"artifacts",
	"bad art",
	// "gross",
	"poor quality",
	"low quality",
}

var enhancedPrompts = []string{
	"8k unity render",
	"skin pores",
	"detailed iris",
	"detailed",
	"detailed face",
	"(vibrant)",
	"sharp focus",
	"(8k)",
}

var HttpClient = http.Client{
	Timeout: 30 * time.Second,
}

// Prompt2PosNeg decomposes a text input into positive and negative prompts
func (s *Server) Prompt2PosNeg(input string) (string, string) {
	left, right, ok := strings.Cut(input, "--")
	if ok {
		return left, right
	}
	return input, s.NegativePrompt
}

var PromptReplace = strings.NewReplacer(
	"self portrait", "painting of an ugly green goblin painting a self portrait on an easel. this is in a basement",
)

func (s *Server) Txt2Img(prompt string) ([]byte, error) {
	// quick hack
	p := NewTxt2ImgParameters()
	p.RestoreFaces = true
	p.Width = 768
	p.Height = 768
	p.Steps = 30

	// if len(prompt) < 20 {
	// 	p.Prompt = p.Prompt + ", " + enhancedPrompt
	// }
	prompt = PromptReplace.Replace(prompt)

	pos, neg := s.Prompt2PosNeg(prompt)
	p.Prompt = pos
	p.NegativePrompt = neg
	p.SetSampler("DPM++ SDE")

	u := s.URL.String()
	u += "/sdapi/v1/txt2img"
	log.Println(u, p.IoReader().String())
	t1 := time.Now()
	resp, err := HttpClient.Post(u, "application/json", p.IoReader())
	if err != nil || resp == nil || resp.Body == nil {
		log.Println(err)
		return nil, err
	}
	t2 := time.Now()
	log.Println("Time taken", t2.Sub(t1))
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
}

// /sdapi/v1/options
func (s *Server) GetConfig() (*Config, error) {
	u := s.URL.String()
	u += "/sdapi/v1/options"
	resp, err := http.Get(u)
	if err != nil || resp == nil || resp.Body == nil {
		log.Println(err)
		return nil, err
	}
	config := &Config{}
	body := resp.Body
	defer body.Close()
	err = json.NewDecoder(body).Decode(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return config, nil
}

// /sdapi/v1/sd-models
func (s *Server) GetModels() ([]Model, error) {
	u := s.URL.String()
	u += "/sdapi/v1/sd-models"
	resp, err := http.Get(u)
	if err != nil || resp == nil || resp.Body == nil {
		log.Println(err)
		return nil, err
	}
	models := []Model{}
	body := resp.Body
	defer body.Close()
	err = json.NewDecoder(body).Decode(&models)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, v := range models {
		s.Models[v.Hash] = v
	}

	return models, nil
}

func (s *Server) SetConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	u := s.URL.String()
	u += "/sdapi/v1/options"
	log.Println(u, config.IoReader().String())
	resp, err := http.Post(u, "application/json", config.IoReader())
	if err != nil || resp == nil || resp.Body == nil {
		log.Println(err)
		return err
	}
	// FIXME: check http errorcode, etc

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	_ = body
	// fmt.Println(string(body))
	return nil
}

func (c *Config) SetModel(hash string, s *Server) {
	if s == nil {
		return
	}
	m, ok := s.Models[hash]
	if !ok {
		return
	}
	c.SdModelCheckpoint = m.Title
}
