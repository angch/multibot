package compreface

import (
	"log"
	"net/url"
)

// Impleements
// https://github.com/exadel-inc/CompreFace/blob/master/docs/Rest-API-description.md#face-detection-service

type Compreface struct {
	BaseURL url.URL
	Options ComprefaceOptions // FIXME
	ApiKey  string

	RecognitionService  *RecognitionService
	VerificationService string
	DetectionService    string
}

type ComprefaceOptions struct {
	DetProbThreshold float64
	Limit            int
	PredictionCount  int
	FacePlugins      string
	Status           bool
}

type RecognitionService struct {
	URL    url.URL
	ApiKey string
}

type FaceCollection struct {
}

type RecognizeResult struct {
	Age struct {
		Probability float64 `json:"probability"`
		High        int     `json:"high"`
		Low         int     `json:"low"`
	} `json:"age"`
	Gender struct {
		Probability float64 `json:"probability"`
		Value       string  `json:"value"`
	} `json:"gender"`
	Mask struct {
		Probability float64 `json:"probability"`
		Value       string  `json:"value"`
	} `json:"mask"`
	// Embedding []interface{} `json:"embedding"`
	Box struct {
		Probability float64 `json:"probability"`
		XMax        int     `json:"x_max"`
		YMax        int     `json:"y_max"`
		XMin        int     `json:"x_min"`
		YMin        int     `json:"y_min"`
	} `json:"box"`
	Landmarks [][]int `json:"landmarks"`
	Subjects  []struct {
		Similarity float64 `json:"similarity"`
		Subject    string  `json:"subject"`
	} `json:"subjects"`
	ExecutionTime struct {
		Age        float64 `json:"age"`
		Gender     float64 `json:"gender"`
		Detector   float64 `json:"detector"`
		Calculator float64 `json:"calculator"`
		Mask       float64 `json:"mask"`
	} `json:"execution_time"`
}

type RecognizeResults struct {
	Result          []RecognizeResult `json:"result"`
	PluginsVersions struct {
		Age        string `json:"age"`
		Gender     string `json:"gender"`
		Detector   string `json:"detector"`
		Calculator string `json:"calculator"`
		Mask       string `json:"mask"`
	} `json:"plugins_versions"`
}

func New(baseURL string) *Compreface {
	return NewWithOptions(baseURL, &ComprefaceOptions{})
}

func NewWithOptions(baseURL string, opt *ComprefaceOptions) *Compreface {
	myUrl, err := url.Parse(baseURL)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &Compreface{
		BaseURL: *myUrl,
		Options: *opt,
	}
}

func (c *Compreface) InitFaceRecognition(apikey string) *RecognitionService {
	if c == nil {
		return nil
	}
	c.RecognitionService = &RecognitionService{
		ApiKey: apikey,
		URL:    c.BaseURL,
	}
	return c.RecognitionService
}

func (c *Compreface) InitFaceVerification(apikey string) {
}

func (c *Compreface) InitFaceDetection(apikey string) {
}

func (r *FaceCollection) Add(imagePath string, subject string) {

}
