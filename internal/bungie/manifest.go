package bungie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BungieBaseURL = "https://www.bungie.net"
	ManifestAPI   = BungieBaseURL + "/Platform/Destiny2/Manifest/"
	UserAgent     = "ContribEmblem/1.0 (+https://github.com/castrojo/contribemblem)"
	ManifestCache = "data/manifest.json"
	EmblemOutput  = "data/emblem.jpg"
)

// Manifest API response structures
type manifestResponse struct {
	ErrorCode   int    `json:"ErrorCode"`
	ErrorStatus string `json:"ErrorStatus"`
	Response    struct {
		JSONWorldComponentContentPaths struct {
			En struct {
				DestinyInventoryItemDefinition string `json:"DestinyInventoryItemDefinition"`
			} `json:"en"`
		} `json:"jsonWorldComponentContentPaths"`
	} `json:"Response"`
}

// Emblem data from manifest
type emblemData struct {
	DisplayProperties struct {
		Icon string `json:"icon"`
	} `json:"displayProperties"`
	SecondaryIcon    string `json:"secondaryIcon"`    // 474x96 wide banner
	SecondarySpecial string `json:"secondarySpecial"` // high-res detail view (1920x1080+)
}

// FetchEmblem downloads emblem artwork from Bungie API
// emblemHash: emblem identifier (e.g., "1409726931")
// Returns path to downloaded emblem image
func FetchEmblem(emblemHash string) error {
	apiKey := os.Getenv("BUNGIE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNGIE_API_KEY environment variable not set")
	}

	fmt.Fprintf(os.Stderr, "Fetching emblem hash: %s\n", emblemHash)

	// Fetch manifest metadata
	fmt.Fprintf(os.Stderr, "Fetching Bungie manifest metadata...\n")
	manifestURL, err := getManifestURL(apiKey)
	if err != nil {
		return fmt.Errorf("failed to get manifest URL: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Manifest URL: %s\n", manifestURL)

	// Download manifest if not cached
	if err := downloadManifestIfNeeded(manifestURL); err != nil {
		return fmt.Errorf("failed to download manifest: %w", err)
	}

	// Look up emblem in manifest
	fmt.Fprintf(os.Stderr, "Looking up emblem %s in manifest...\n", emblemHash)
	iconPath, err := lookupEmblemIcon(emblemHash)
	if err != nil {
		return fmt.Errorf("failed to lookup emblem: %w", err)
	}

	// Download emblem image
	iconURL := BungieBaseURL + iconPath
	fmt.Fprintf(os.Stderr, "Downloading emblem image from: %s\n", iconURL)
	if err := downloadImage(iconURL, EmblemOutput); err != nil {
		return fmt.Errorf("failed to download emblem image: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Emblem image saved to %s\n", EmblemOutput)
	return nil
}

func getManifestURL(apiKey string) (string, error) {
	req, err := http.NewRequest("GET", ManifestAPI, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Log rate limit
	if remaining := resp.Header.Get("X-Ratelimit-Remaining"); remaining != "" {
		fmt.Fprintf(os.Stderr, "ℹ️  Rate limit remaining: %s\n", remaining)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var manifest manifestResponse
	if err := json.Unmarshal(body, &manifest); err != nil {
		return "", err
	}

	if manifest.ErrorCode != 1 {
		return "", fmt.Errorf("Bungie API error %d: %s", manifest.ErrorCode, manifest.ErrorStatus)
	}

	manifestURL := manifest.Response.JSONWorldComponentContentPaths.En.DestinyInventoryItemDefinition
	if manifestURL == "" {
		return "", fmt.Errorf("manifest URL not found in response")
	}

	return BungieBaseURL + manifestURL, nil
}

func downloadManifestIfNeeded(url string) error {
	// Check if manifest exists and is fresh (< 24 hours old)
	if info, err := os.Stat(ManifestCache); err == nil {
		age := time.Since(info.ModTime())
		if age < 24*time.Hour {
			fmt.Fprintf(os.Stderr, "✓ Using cached manifest (age: %v)\n", age.Round(time.Minute))
			return nil
		}
		fmt.Fprintf(os.Stderr, "⚠️  Manifest cache expired (age: %v), re-downloading...\n", age.Round(time.Minute))
	}

	fmt.Fprintf(os.Stderr, "Downloading manifest database (~100MB, this may take a moment)...\n")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	out, err := os.Create(ManifestCache)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "✓ Manifest cached\n")
	return nil
}

func lookupEmblemIcon(emblemHash string) (string, error) {
	data, err := os.ReadFile(ManifestCache)
	if err != nil {
		return "", err
	}

	var manifest map[string]emblemData
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}

	emblem, ok := manifest[emblemHash]
	if !ok {
		return "", fmt.Errorf("emblem hash %s not found in manifest", emblemHash)
	}

	// Prefer secondarySpecial (high-res 1920x1080+) over secondaryIcon (474x96)
	// This allows downscaling instead of upscaling for sharper results
	if emblem.SecondarySpecial != "" {
		return emblem.SecondarySpecial, nil
	}

	// Fall back to secondaryIcon (474x96 wide banner) if secondarySpecial not available
	if emblem.SecondaryIcon != "" {
		return emblem.SecondaryIcon, nil
	}

	if emblem.DisplayProperties.Icon == "" {
		return "", fmt.Errorf("icon path not found for emblem %s", emblemHash)
	}

	return emblem.DisplayProperties.Icon, nil
}

func downloadImage(url, outputPath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	// Save raw download directly to avoid JPEG re-encoding artifacts
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
