package sdapi

type HTTPValidationError struct {
	Detail []Detail `json:"detail"`
}

type Detail struct {
	Loc  []any  `json:"loc"`
	Msg  string `json:"msg"`
	Type string `json:"type"`
}
