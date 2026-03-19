package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrGeminiNotConfigured = errors.New("gemini api key not configured")

type AIContentService struct {
	apiKey     string
	httpClient *http.Client
}

func NewAIContentService(apiKey string) *AIContentService {
	return &AIContentService{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

type GenerateContentInput struct {
	Text   string   `json:"text"`
	Formats []string `json:"formats"` // "tweet", "reel_script", "hook"
}

type GenerateContentOutput struct {
	Tweet      string `json:"tweet,omitempty"`
	ReelScript string `json:"reel_script,omitempty"`
	Hook       string `json:"hook,omitempty"`
}

func (s *AIContentService) Generate(ctx context.Context, input GenerateContentInput) (*GenerateContentOutput, error) {
	if s.apiKey == "" {
		return nil, ErrGeminiNotConfigured
	}
	prompt := buildPrompt(input)
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]interface{}{{"text": prompt}}},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 1024,
			"temperature":     0.7,
		},
	}
	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", s.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini api error: %d", resp.StatusCode)
	}
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("no content in response")
	}
	text := result.Candidates[0].Content.Parts[0].Text
	return parseGeneratedText(text, input.Formats), nil
}

func buildPrompt(input GenerateContentInput) string {
	formats := input.Formats
	if len(formats) == 0 {
		formats = []string{"tweet", "reel_script", "hook"}
	}
	var want []string
	for _, f := range formats {
		switch f {
		case "tweet":
			want = append(want, "a short tweet (under 280 chars)")
		case "reel_script":
			want = append(want, "a short reel script (3-5 lines)")
		case "hook":
			want = append(want, "a punchy hook (one line)")
		}
	}
	if len(want) == 0 {
		want = []string{"a tweet", "a reel script", "a hook"}
	}
	return fmt.Sprintf("Based on this input, generate the following. Keep each section concise.\n\nInput:\n%s\n\nGenerate:\n- %s\n\nRespond with clear labels like TWEET:, REEL SCRIPT:, HOOK: so the response can be parsed.", input.Text, strings.Join(want, "\n- "))
}

func parseGeneratedText(text string, _ []string) *GenerateContentOutput {
	out := &GenerateContentOutput{}
	lines := strings.Split(text, "\n")
	var currentKey string
	var current []string
	flush := func() {
		s := strings.TrimSpace(strings.Join(current, " "))
		switch currentKey {
		case "tweet":
			out.Tweet = s
		case "reel_script":
			out.ReelScript = s
		case "hook":
			out.Hook = s
		}
		current = nil
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		upper := strings.ToUpper(line)
		if strings.HasPrefix(upper, "TWEET:") {
			flush()
			currentKey = "tweet"
			current = []string{strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "TWEET:"), "tweet:"))}
		} else if strings.HasPrefix(upper, "REEL SCRIPT:") {
			flush()
			currentKey = "reel_script"
			current = []string{strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "REEL SCRIPT:"), "reel script:"))}
		} else if strings.HasPrefix(upper, "HOOK:") {
			flush()
			currentKey = "hook"
			current = []string{strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "HOOK:"), "hook:"))}
		} else if currentKey != "" && line != "" {
			current = append(current, line)
		}
	}
	flush()
	if out.Tweet == "" && out.ReelScript == "" && out.Hook == "" {
		out.Tweet = strings.TrimSpace(text)
	}
	return out
}
