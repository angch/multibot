package sdapi

type Model struct {
	Title     string `json:"title"`
	ModelName string `json:"model_name"`
	Hash      string `json:"hash"`
	Filename  string `json:"filename"`
	Config    string `json:"config"`
}
