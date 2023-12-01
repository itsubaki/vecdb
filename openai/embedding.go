package openai

type Request struct {
	Input          []string `json:"input"`
	Model          string   `json:"model"`
	EncodingFormat string   `json:"encoding_format"`
}

type Response struct {
	Object string `json:"object"`
	Model  string `json:"model"`
	Data   []Data `json:"data"`
	Usage  Usage  `json:"usage"`
}

type Data struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding Embedding `json:"embedding"`
}

type Embedding []float64

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
