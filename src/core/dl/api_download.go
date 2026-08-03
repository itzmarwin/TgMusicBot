/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package dl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"ashokshau/tgmusic/config"
)

const (
	shrutiApiUrl = "https://api.shrutibots.site"
	shrutiApiKey = ""
)

type apiData struct {
	VideoID string
}

func newApiDownload(videoID string) *apiData {
	return &apiData{VideoID: videoID}
}

func mediaType(video bool) string {
	if video {
		return "video"
	}
	return "audio"
}

func fileExt(video bool) string {
	if video {
		return "mp4"
	}
	return "mp3"
}

func (a *apiData) Process(video bool) (string, error) {
	if a.VideoID == "" {
		return "", errors.New("videoID is empty")
	}

	if err := os.MkdirAll(config.DownloadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create downloads dir: %w", err)
	}

	filePath := filepath.Join(config.DownloadsDir, fmt.Sprintf("%s.%s", a.VideoID, fileExt(video)))

	if info, err := os.Stat(filePath); err == nil && info.Size() > 0 {
		return filePath, nil
	}

	timeout := 5 * time.Minute
	if video {
		timeout = 10 * time.Minute
	}
	mType := mediaType(video)

	if config.DevilApiUrl != "" && config.DevilApiKey != "" {
		if err := fetchToFile(config.DevilApiUrl, config.DevilApiKey, a.VideoID, mType, filePath, timeout); err == nil {
			slog.Info("downloaded via devil api", "video_id", a.VideoID)
			return filePath, nil
		} else {
			slog.Warn("devil api failed, falling back to shruti", "video_id", a.VideoID, "error", err)
		}
	}

	if shrutiApiUrl == "" || shrutiApiKey == "" {
		return "", errors.New("no download API is configured")
	}

	if err := fetchToFile(shrutiApiUrl, shrutiApiKey, a.VideoID, mType, filePath, timeout); err != nil {
		_ = os.Remove(filePath)
		return "", fmt.Errorf("shruti api failed: %w", err)
	}

	slog.Info("downloaded via shruti api", "video_id", a.VideoID)
	return filePath, nil
}

func fetchToFile(baseURL, apiKey, videoID, mType, destPath string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reqURL, err := url.Parse(baseURL + "/download")
	if err != nil {
		return fmt.Errorf("invalid api url: %w", err)
	}

	q := reqURL.Query()
	q.Set("url", videoID)
	q.Set("type", mType)
	q.Set("api_key", apiKey)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	partPath := destPath + ".part"
	out, err := os.Create(partPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	const chunkSize = 131072
	buf := make([]byte, chunkSize)
	_, copyErr := io.CopyBuffer(out, resp.Body, buf)
	closeErr := out.Close()

	if copyErr != nil {
		_ = os.Remove(partPath)
		return fmt.Errorf("failed to write file: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(partPath)
		return fmt.Errorf("failed to close file: %w", closeErr)
	}

	if info, err := os.Stat(partPath); err != nil || info.Size() == 0 {
		_ = os.Remove(partPath)
		return errors.New("downloaded file is empty")
	}

	if err := os.Rename(partPath, destPath); err != nil {
		return fmt.Errorf("failed to finalize file: %w", err)
	}

	return nil
}
