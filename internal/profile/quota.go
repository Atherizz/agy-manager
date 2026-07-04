package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type loadCodeAssistResponse struct {
	CloudAiCompanionProject string `json:"cloudaicompanionProject"`
}

type quotaApiResponse struct {
	Models map[string]struct {
		QuotaInfo struct {
			RemainingFraction float64 `json:"remainingFraction"`
		} `json:"quotaInfo"`
	} `json:"models"`
}

// GetQuota fetches the remaining quota fraction for current models
func GetQuota(profileName, activeProfile string) (map[string]float64, error) {
	token := GetAccessToken(profileName, activeProfile)
	if token == "" {
		return nil, fmt.Errorf("no access token found (profile might not be logged in)")
	}

	client := &http.Client{Timeout: 15 * time.Second}

	// 1. Fetch Project ID
	req1, _ := http.NewRequest("POST", "https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist", bytes.NewBuffer([]byte(`{"metadata": {"ideType": "ANTIGRAVITY"}}`)))
	req1.Header.Set("Authorization", "Bearer "+token)
	req1.Header.Set("User-Agent", "antigravity/windows/amd64")
	req1.Header.Set("Content-Type", "application/json")
	
	res1, err := client.Do(req1)
	if err != nil {
		return nil, fmt.Errorf("failed to load code assist: %w", err)
	}
	defer res1.Body.Close()

	var pData loadCodeAssistResponse
	json.NewDecoder(res1.Body).Decode(&pData)

	// 2. Fetch Models Quota
	payload := []byte(`{}`)
	if pData.CloudAiCompanionProject != "" {
		payload = []byte(fmt.Sprintf(`{"project": "%s"}`, pData.CloudAiCompanionProject))
	}

	req2, _ := http.NewRequest("POST", "https://cloudcode-pa.googleapis.com/v1internal:fetchAvailableModels", bytes.NewBuffer(payload))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("User-Agent", "antigravity/1.11.3 Darwin/arm64")
	req2.Header.Set("Content-Type", "application/json")

	res2, err := client.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}
	defer res2.Body.Close()

	if res2.StatusCode != 200 {
		body, _ := io.ReadAll(res2.Body)
		return nil, fmt.Errorf("API error %d: %s", res2.StatusCode, string(body))
	}

	var quotaRes quotaApiResponse
	if err := json.NewDecoder(res2.Body).Decode(&quotaRes); err != nil {
		return nil, err
	}

	results := make(map[string]float64)
	for modelName, info := range quotaRes.Models {
		// Filter out old gemini-1.x and 2.x models
		if strings.Contains(modelName, "gemini-1") || strings.Contains(modelName, "gemini-2") {
			continue
		}
		if strings.Contains(modelName, "gemini") || strings.Contains(modelName, "claude") {
			results[modelName] = info.QuotaInfo.RemainingFraction
		}
	}

	return results, nil
}
