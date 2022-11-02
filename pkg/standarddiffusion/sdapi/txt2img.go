package sdapi

import (
	"bytes"
	"encoding/json"
)

type Info struct {
	Prompt                string      `json:"prompt"`
	AllPrompts            []string    `json:"all_prompts"`
	NegativePrompt        string      `json:"negative_prompt"`
	Seed                  int         `json:"seed"`
	AllSeeds              []int       `json:"all_seeds"`
	Subseed               int         `json:"subseed"`
	AllSubseeds           []int       `json:"all_subseeds"`
	SubseedStrength       float64     `json:"subseed_strength"`
	Width                 int         `json:"width"`
	Height                int         `json:"height"`
	SamplerIndex          int         `json:"sampler_index"`
	Sampler               string      `json:"sampler"`
	CfgScale              float64     `json:"cfg_scale"`
	Steps                 int         `json:"steps"`
	BatchSize             int         `json:"batch_size"`
	RestoreFaces          bool        `json:"restore_faces"`
	FaceRestorationModel  interface{} `json:"face_restoration_model"`
	SdModelHash           string      `json:"sd_model_hash"`
	SeedResizeFromW       int         `json:"seed_resize_from_w"`
	SeedResizeFromH       int         `json:"seed_resize_from_h"`
	DenoisingStrength     float64     `json:"denoising_strength"`
	ExtraGenerationParams struct {
	} `json:"extra_generation_params"`
	IndexOfFirstImage int      `json:"index_of_first_image"`
	Infotexts         []string `json:"infotexts"`
	Styles            []string `json:"styles"`
	JobTimestamp      string   `json:"job_timestamp"`
	ClipSkip          int      `json:"clip_skip"`
}

type Txt2ImgParameters struct {
	EnableHr          bool     `json:"enable_hr"`
	DenoisingStrength float64  `json:"denoising_strength"`
	FirstphaseWidth   int      `json:"firstphase_width"`
	FirstphaseHeight  int      `json:"firstphase_height"`
	Prompt            string   `json:"prompt"`
	Styles            []string `json:"styles"`
	Seed              int      `json:"seed"`
	Subseed           int      `json:"subseed"`
	SubseedStrength   float64  `json:"subseed_strength"`
	SeedResizeFromH   int      `json:"seed_resize_from_h"`
	SeedResizeFromW   int      `json:"seed_resize_from_w"`
	BatchSize         int      `json:"batch_size"`
	NIter             int      `json:"n_iter"`
	Steps             int      `json:"steps"`
	CfgScale          float64  `json:"cfg_scale"`
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	RestoreFaces      bool     `json:"restore_faces"`
	Tiling            bool     `json:"tiling"`
	NegativePrompt    string   `json:"negative_prompt"`
	Eta               float64  `json:"eta"`
	SChurn            float64  `json:"s_churn"`
	STmax             float64  `json:"s_tmax"`
	STmin             float64  `json:"s_tmin"`
	SNoise            float64  `json:"s_noise"`
	OverrideSettings  struct {
	} `json:"override_settings"`
	SamplerIndex string `json:"sampler_index"`
}

type Txt2ImgParametersResult struct {
	Images     []string          `json:"images"`
	Parameters Txt2ImgParameters `json:"parameters"`
	Info       string            `json:"info"`
}

func NewTxt2ImgParameters() *Txt2ImgParameters {
	s := Txt2ImgParameters{}
	s.Styles = []string{}
	s.Seed = -1
	s.Subseed = -1
	s.SeedResizeFromH = -1
	s.SeedResizeFromW = -1
	s.BatchSize = 1
	s.NIter = 1
	s.Steps = 50
	s.CfgScale = 7
	s.Width = 512
	s.Height = 512
	s.SNoise = -1
	s.SamplerIndex = "Euler"
	return &s
}

func (p *Txt2ImgParameters) IoReader() *bytes.Buffer {
	j, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return bytes.NewBuffer(j)
}
