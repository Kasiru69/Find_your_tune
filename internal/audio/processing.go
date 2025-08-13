package audio

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"Shazam/internal/database"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	SpotifyClientID     = "eec03041bad34931a01c2d8106bef880"
	SpotifyClientSecret = "66ea4b4480034839ae27ab41a9a20d1b"
)

func ProcessAudioFile(filePath string) (*AudioFingerprint, error) {
	cmd := exec.Command("ffmpeg", "-i", filePath, "-f", "f64le", "-acodec", "pcm_f64le", "-ac", "1", "-ar", "22050", "-")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to process audio: %v", err)
	}

	samples := make([]float64, len(output)/8)
	for i := 0; i < len(samples); i++ {
		bits := uint64(0)
		for j := 0; j < 8; j++ {
			bits |= uint64(output[i*8+j]) << (j * 8)
		}
		samples[i] = math.Float64frombits(bits)
	}

	return GenerateFingerprint(samples)
}

func DownloadAudioPreview(searchQuery string, outputPath string) error {
	fmt.Printf("Downloading: %s\n", searchQuery)
	os.Remove(outputPath)

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "wav",
		"--audio-quality", "0",
		"--default-search", "ytsearch1:",
		"--no-warnings",
		"-o", outputPath,
		searchQuery)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("download failed: %v\nOutput: %s", err, string(output))
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("audio file was not created")
	}

	fmt.Println("Download successful!")
	return nil
}

func AddSongToDatabase(db *database.DB, artistName, songName, albumName string) error {
	fmt.Printf("🎵 Adding song to database: %s - %s\n", artistName, songName)

	searchQuery := fmt.Sprintf("%s %s", artistName, songName)
	tempFile := fmt.Sprintf("data/temp/temp_%d.wav", time.Now().Unix())

	err := DownloadAudioPreview(searchQuery, tempFile)
	if err != nil {
		return fmt.Errorf("failed to download %s - %s: %v", artistName, songName, err)
	}
	defer os.Remove(tempFile)

	fingerprint, err := ProcessAudioFile(tempFile)
	if err != nil {
		return fmt.Errorf("failed to generate fingerprint: %v", err)
	}

	song := &database.Song{
		Title:        songName,
		Artist:       artistName,
		Album:        albumName,
		Duration:     len(fingerprint.HashSegments) * 512 / 22050,
		Fingerprint:  fingerprint.Fingerprint,
		HashSegments: fingerprint.HashSegments,
	}

	return db.AddSong(song)
}

func ExtractSpotifyID(input string) (string, error) {
	u := strings.TrimSpace(input)
	if u == "" {
		return "", fmt.Errorf("empty input")
	}

	if strings.HasPrefix(u, "spotify:track:") {
		id := strings.TrimPrefix(u, "spotify:track:")
		if len(id) >= 21 && len(id) <= 22 {
			return id, nil
		}
	}

	if parsed, err := url.Parse(u); err == nil {
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		for i := 0; i < len(parts); i++ {
			if parts[i] == "track" && i+1 < len(parts) {
				id := parts[i+1]
				id = strings.SplitN(id, "?", 2)[0]
				if len(id) >= 21 && len(id) <= 22 {
					return id, nil
				}
				break
			}
		}
	}

	re := regexp.MustCompile(`([A-Za-z0-9]{21,22})`)
	if m := re.FindStringSubmatch(u); m != nil {
		return m[1], nil
	}

	return "", fmt.Errorf("invalid Spotify URL format")
}

func GetTrackInfo(trackID string) (*spotify.FullTrack, error) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     SpotifyClientID,
		ClientSecret: SpotifyClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get token: %v", err)
	}
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	track, err := client.GetTrack(ctx, spotify.ID(trackID))
	if err != nil {
		return nil, fmt.Errorf("couldn't get track: %v", err)
	}
	return track, nil
}
