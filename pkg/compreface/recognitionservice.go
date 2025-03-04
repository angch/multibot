package compreface

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

func (r *RecognitionService) Call(path string) []byte {
	if r == nil {
		return nil
	}
	myUrl := r.URL
	myUrl.Path = path
	u := myUrl.String()
	data, err := getapi(u, r.ApiKey)
	if err != nil {
		return nil
	}
	return data
}

func (r *RecognitionService) Post(path string, params url.Values, postdata []byte) []byte {
	if r == nil {
		return nil
	}
	myUrl := r.URL
	myUrl.Path = path
	myUrl.RawQuery = params.Encode()
	u := myUrl.String()
	data, err := postapi(u, r.ApiKey, postdata)
	if err != nil {
		return nil
	}
	return data
}

func (r *RecognitionService) GetFaceCollection() *FaceCollection {
	return &FaceCollection{}
}

type GetSubjectsResult struct {
	Subjects []string `json:"subjects"`
}

func (r *RecognitionService) GetSubjects() []string {
	b := r.Call("/api/v1/recognition/subjects")
	subjects := GetSubjectsResult{}
	err := json.Unmarshal(b, &subjects)
	if err != nil {
		log.Println(err)
		return nil
	}
	return subjects.Subjects
}

func (r *RecognitionService) Recognize(imagePath string) {

}

func (r *RecognitionService) AddFace(subject string, det_prob_threshold float64, data []byte) []byte {
	params := url.Values{
		"subject":            []string{subject},
		"det_prob_threshold": []string{fmt.Sprintf("%.2f", det_prob_threshold)},
	}
	b := r.Post("/api/v1/recognition/faces", params, data)
	return b
}
