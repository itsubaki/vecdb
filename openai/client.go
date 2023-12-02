package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const V1 = "https://api.openai.com/v1"

type Client struct {
	Org     string
	APIKey  string
	ModelID string
}

func (c *Client) Models() (*Models, error) {
	url, err := url.JoinPath(V1, "/models")
	if err != nil {
		return nil, fmt.Errorf("join path: %v", err)
	}

	resp, err := c.do("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg Error
		if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
			return nil, fmt.Errorf("decode: %v", err)
		}

		return nil, fmt.Errorf("status code=%v, message: %v", resp.StatusCode, msg.Error.Message)
	}

	var models Models
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return &models, nil
}

func (c *Client) Embeddings(text []string) ([][]float64, error) {
	url, err := url.JoinPath(V1, "/embeddings")
	if err != nil {
		return nil, fmt.Errorf("join path: %v", err)
	}

	req := Request{
		Model:          c.ModelID,
		Input:          text,
		EncodingFormat: "float",
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal: %v", err)
	}

	resp, err := c.do("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}
	defer resp.Body.Close()

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	out := make([][]float64, len(res.Data))
	for i, d := range res.Data {
		out[i] = d.Embedding
	}

	return out, nil
}

func (c *Client) do(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Organization", c.Org)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.APIKey))

	return http.DefaultClient.Do(req)
}
