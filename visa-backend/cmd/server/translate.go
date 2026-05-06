package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type translateTextRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	Mode       string `json:"mode"`
}

type translateTextResponse struct {
	TranslatedText string `json:"translatedText"`
	SourceLang     string `json:"sourceLang"`
	TargetLang     string `json:"targetLang"`
	Provider       string `json:"provider"`
}

type openAIChatCompletionRequest struct {
	Model       string              `json:"model"`
	Messages    []openAIChatMessage `json:"messages"`
	Temperature float64             `json:"temperature,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatCompletionResponse struct {
	Choices []struct {
		Message openAIChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *appServer) handleTranslateText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	var in translateTextRequest
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	in.Text = strings.TrimSpace(in.Text)
	in.SourceLang = normalizeTranslateLang(in.SourceLang)
	in.TargetLang = normalizeTranslateLang(in.TargetLang)
	in.Mode = strings.TrimSpace(in.Mode)

	if in.Text == "" {
		writeError(w, http.StatusBadRequest, errors.New("text is required"))
		return
	}
	if len([]rune(in.Text)) > 2000 {
		writeError(w, http.StatusBadRequest, errors.New("text is too long"))
		return
	}
	if !isTranslateLangSupported(in.SourceLang) || !isTranslateLangSupported(in.TargetLang) {
		writeError(w, http.StatusBadRequest, errors.New("only zh and en are supported"))
		return
	}
	if in.SourceLang == in.TargetLang {
		writeJSON(w, http.StatusOK, apiResponse{
			Message: "ok",
			Data: translateTextResponse{
				TranslatedText: in.Text,
				SourceLang:     in.SourceLang,
				TargetLang:     in.TargetLang,
				Provider:       "local",
			},
		})
		return
	}

	out, err := s.translateText(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: out})
}

func normalizeTranslateLang(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch lang {
	case "zh-cn", "zh-hans", "cn":
		return "zh"
	case "en-us", "en-gb":
		return "en"
	default:
		return lang
	}
}

func isTranslateLangSupported(lang string) bool {
	switch lang {
	case "zh", "en":
		return true
	default:
		return false
	}
}

func (s *appServer) translateText(ctx context.Context, in translateTextRequest) (translateTextResponse, error) {
	if strings.TrimSpace(s.cfg.OpenAIAPIKey) == "" {
		return translateTextResponse{}, errors.New("translation service is not configured")
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	systemPrompt := buildTranslateSystemPrompt(in.SourceLang, in.TargetLang, in.Mode)
	reqBody := openAIChatCompletionRequest{
		Model: strings.TrimSpace(s.cfg.OpenAITranslateModel),
		Messages: []openAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: in.Text},
		},
		Temperature: 0.2,
		MaxTokens:   1200,
	}
	if reqBody.Model == "" {
		reqBody.Model = "gpt-4o-mini"
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return translateTextResponse{}, err
	}

	baseURL := strings.TrimRight(strings.TrimSpace(s.cfg.OpenAIBaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return translateTextResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(s.cfg.OpenAIAPIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return translateTextResponse{}, fmt.Errorf("translation request failed: %w", err)
	}
	defer resp.Body.Close()

	var out openAIChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return translateTextResponse{}, fmt.Errorf("failed to parse translation response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if out.Error != nil && strings.TrimSpace(out.Error.Message) != "" {
			return translateTextResponse{}, errors.New(strings.TrimSpace(out.Error.Message))
		}
		return translateTextResponse{}, fmt.Errorf("translation service returned status %d", resp.StatusCode)
	}
	if len(out.Choices) == 0 {
		return translateTextResponse{}, errors.New("translation result is empty")
	}

	translated := strings.TrimSpace(out.Choices[0].Message.Content)
	if translated == "" {
		return translateTextResponse{}, errors.New("translation result is empty")
	}

	return translateTextResponse{
		TranslatedText: translated,
		SourceLang:     in.SourceLang,
		TargetLang:     in.TargetLang,
		Provider:       "openai-compatible",
	}, nil
}

func buildTranslateSystemPrompt(sourceLang, targetLang, mode string) string {
	sourceLabel := translateLangLabel(sourceLang)
	targetLabel := translateLangLabel(targetLang)

	base := []string{
		"You are a professional translation engine for visa and travel preparation content.",
		fmt.Sprintf("Translate the user's text from %s to %s.", sourceLabel, targetLabel),
		"Return only the translated text.",
		"Do not explain, annotate, or add commentary.",
		"Preserve names, passport numbers, dates, addresses, phone numbers, amounts, codes, and formatting unless translation is required.",
		"If the input already contains mixed bilingual content, keep the parts already in the target language natural and only translate the rest.",
	}
	if strings.EqualFold(strings.TrimSpace(mode), "visa") {
		base = append(base,
			"Use a formal and practical tone suitable for visa materials, travel documents, application notes, itineraries, and supporting letters.",
			"Prefer accurate and concise wording over literary style.",
		)
	}
	return strings.Join(base, "\n")
}

func translateLangLabel(lang string) string {
	switch lang {
	case "zh":
		return "Simplified Chinese"
	case "en":
		return "English"
	default:
		return strings.ToUpper(lang)
	}
}
