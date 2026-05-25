package httpx

type Problem struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}
