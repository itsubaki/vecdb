package openai

const (
	GPT35_TURBO            string = "gpt-3.5-turbo"
	GPT4                   string = "gpt-4"
	TEXT_EMBEDDING_ADA_002 string = "text-embedding-ada-002"
)

type Models struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}
