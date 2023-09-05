package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// JSON request
type RequestBody struct {
	Instances  []Instance `json:"instances"`
	Parameters Parameters `json:"parameters"`
}

type Instance struct {
	Context  string    `json:"context,omitempty"`
	Examples []Example `json:"examples,omitempty"`
	Messages []Message `json:"messages"`
}

type Example struct {
	Input struct {
		Content string `json:"content"`
	} `json:"input"`
	Output struct {
		Content string `json:"content"`
	} `json:"output"`
}

type Message struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type Parameters struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int64   `json:"maxOutputTokens"`
	TopP            float64 `json:"topP"`
	TopK            int64   `json:"topK"`
}

// JSON response
type ResponseBody struct {
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Candidates       []Candidate `json:"candidates"`
	CitationMetadata []struct {
		Citations []Citation `json:"citations"`
	} `json:"citationMetadata"`
	SafetyAttributes []SafetyAttribute `json:"safetyAttributes"`
	Metadata         struct {
		TokenMetadata `json:"tokenMetadata"`
	} `json:"metadata"`
}

type Candidate struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type Citation struct {
	StartIndex      int64  `json:"startIndex"`
	EndIndex        int64  `json:"endIndex"`
	URL             string `json:"url"`
	Title           string `json:"title"`
	License         string `json:"license"`
	PublicationDate string `json:"publicationDate"`
}

type SafetyAttribute struct {
	Categories []string  `json:"categories"`
	Blocked    bool      `json:"blocked"`
	Scores     []float64 `json:"scores"`
}

type TokenMetadata struct {
	InputTokenCount  TokenCount `json:"inputTokenCount"`
	OutputTokenCount TokenCount `json:"outputTokenCount"`
}

type TokenCount struct {
	TotalBillableCharacters int64 `json:"totalBillableCharacters"`
	TotalTokens             int64 `json:"totalTokens"`
}

// Implementation
type GCPVertex struct {
	apiEndpoint string
	projectID   string
	modelID     string
	config      *jwt.Config
}

func NewGCPVertex(credentialsPath, apiEndpoint, projectID, modelID string) *GCPVertex {
	credentials, _ := os.ReadFile(credentialsPath)
	config, _ := google.JWTConfigFromJSON(credentials, "https://www.googleapis.com/auth/cloud-platform")
	return &GCPVertex{
		apiEndpoint: apiEndpoint,
		projectID:   projectID,
		modelID:     modelID,
		config:      config,
	}
}

func (l *GCPVertex) GenerateText(ctx context.Context, prompt string, corpus string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v1/projects/%s/locations/us-central1/publishers/google/models/%s:predict",
		l.apiEndpoint,
		l.projectID,
		l.modelID,
	)
	body, err := json.Marshal(RequestBody{
		Instances: []Instance{
			{
				Context: prompt,
				Messages: []Message{
					{
						Author:  "user",
						Content: corpus,
					},
				},
			},
		},
		Parameters: Parameters{
			Temperature:     0.7,
			MaxOutputTokens: 100,
			TopP:            0.8,
			TopK:            40,
		},
	})
	if err != nil {
		return "", err
	}
	client := l.config.Client(ctx)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API call with error: %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("json: %v", string(respBody))
	var response ResponseBody
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}
	return response.Predictions[0].Candidates[0].Content, nil
}
