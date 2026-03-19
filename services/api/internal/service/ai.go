package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var ErrGeminiNotConfigured = errors.New("gemini api key not configured")

var AllowedModels = []string{
	"gemini-2.5-flash",
	"gemini-2.5-flash-lite",
	"gemini-2.5-pro",
}

const DefaultModel = "gemini-2.5-flash"

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

func isValidModel(m string) bool {
	for _, allowed := range AllowedModels {
		if m == allowed {
			return true
		}
	}
	return false
}

type GenerateContentInput struct {
	Text    string   `json:"text"`
	Formats []string `json:"formats"` // "tweet", "reel_script", "hook"
	Model   string   `json:"model"`   // optional: gemini-2.5-flash, gemini-2.5-flash-lite, gemini-2.5-pro
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
	model := input.Model
	if !isValidModel(model) {
		model = DefaultModel
	}
	systemPrompt := buildSystemPrompt(input.Formats)
	userPrompt := buildUserPrompt(input.Text, input.Formats)
	reqBody := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]interface{}{{"text": systemPrompt}},
		},
		"contents": []map[string]interface{}{
			{"parts": []map[string]interface{}{{"text": userPrompt}}},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 1024,
			"temperature":    0.6,
		},
	}
	body, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, s.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody []byte
		errBody, _ = io.ReadAll(resp.Body)
		errMsg := string(errBody)
		if errMsg != "" && len(errMsg) < 500 {
			return nil, fmt.Errorf("gemini api %d: %s", resp.StatusCode, errMsg)
		}
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
	if len(result.Candidates) == 0 {
		return nil, errors.New("gemini returned no candidates (possible safety block or model error)")
	}
	if len(result.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("gemini returned empty content")
	}
	text := result.Candidates[0].Content.Parts[0].Text
	return parseGeneratedText(text, input.Formats), nil
}

func buildSystemPrompt(formats []string) string {
	if len(formats) == 0 {
		formats = []string{"tweet", "reel_script", "hook"}
	}
	return `You are a dev creator content assistant. Turn the user's input (coding notes, LeetCode solutions, ideas) into ready-to-post content.

Rules:
1. Output ONLY the requested formats. Use exactly these labels on their own line: TWEET:, REEL SCRIPT:, HOOK:
2. Tweet: under 280 characters, engaging, 1-2 relevant hashtags (e.g. #Coding #LeetCode).
3. Hook: one punchy line only. No hashtags.
4. Reel script: 3-5 numbered steps. Each step: number, then (brief visual/action in parentheses), then the line in quotes. Example: 1. (Confused look) "Got you stuck?" 2. (Lightbulb) "Here's the fix:"
5. Be concise and actionable. No filler.`
}

func buildUserPrompt(text string, formats []string) string {
	if len(formats) == 0 {
		formats = []string{"tweet", "reel_script", "hook"}
	}
	var want []string
	for _, f := range formats {
		switch f {
		case "tweet":
			want = append(want, "TWEET")
		case "reel_script":
			want = append(want, "REEL SCRIPT")
		case "hook":
			want = append(want, "HOOK")
		}
	}
	if len(want) == 0 {
		want = []string{"TWEET", "REEL SCRIPT", "HOOK"}
	}
	return fmt.Sprintf("Generate only: %s\n\nUser input:\n%s", strings.Join(want, ", "), text)
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
