package openai

type Error struct {
	Error Message `json:"error"`
}

type Message struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}
