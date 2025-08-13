package audio

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func RecordScreenWithAudio(outputFile string, durationSeconds int) error {
	streamPort := "1935"
	httpPort := "8000"

	fmt.Println("🎬 Starting screen recording with audio capture...")

	httpServer := startHTTPServer(httpPort)
	defer httpServer.Process.Kill()

	screenShare := startScreenShare(streamPort)
	defer screenShare.Process.Kill()

	fmt.Printf("Screen sharing started at http://localhost:%s\n", httpPort)
	fmt.Printf("RTMP stream at rtmp://localhost:%s/live/stream\n", streamPort)

	time.Sleep(3 * time.Second)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ffmpeg",
			"-f", "gdigrab",
			"-i", "desktop",
			"-f", "dshow",
			"-i", "audio=Stereo Mix (2- Realtek(R) Audio)",
			"-t", strconv.Itoa(durationSeconds),
			"-vcodec", "libx264",
			"-acodec", "aac",
			"-y", outputFile)
	case "darwin":
		cmd = exec.Command("ffmpeg",
			"-f", "avfoundation",
			"-i", "1:0",
			"-t", strconv.Itoa(durationSeconds),
			"-vcodec", "libx264",
			"-acodec", "aac",
			"-y", outputFile)
	case "linux":
		cmd = exec.Command("ffmpeg",
			"-f", "x11grab",
			"-i", ":0.0",
			"-f", "pulse",
			"-i", "default",
			"-t", strconv.Itoa(durationSeconds),
			"-vcodec", "libx264",
			"-acodec", "aac",
			"-y", outputFile)
	}

	fmt.Printf("🔴 Recording screen + audio for %d seconds...\n", durationSeconds)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("screen recording failed: %v", err)
	}

	fmt.Println("✅ Screen recording with audio completed!")
	return nil
}

func ExtractAudioFingerprint(videoFile string) (*AudioFingerprint, error) {
	fmt.Println("🎵 Extracting audio from screen recording...")

	audioWavFile := "data/temp/recording.wav"

	if _, err := os.Stat(audioWavFile); err == nil {
		os.Remove(audioWavFile)
		fmt.Println("🗑️ Deleted old recording.wav")
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("💾 Creating new recording.wav...")

	saveCmd := exec.Command("ffmpeg",
		"-i", videoFile,
		"-vn",
		"-af", "loudnorm=I=-16:LRA=11:TP=-1.5",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		"-y", audioWavFile)

	output, err := saveCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ FFmpeg output: %s\n", string(output))
		return nil, fmt.Errorf("failed to create new recording.wav: %v", err)
	}

	if fileInfo, err := os.Stat(audioWavFile); err != nil {
		return nil, fmt.Errorf("recording.wav was not created: %v", err)
	} else {
		fmt.Printf("✅ Created fresh recording.wav (%d bytes)\n", fileInfo.Size())
	}

	fmt.Println("🔍 Extracting audio for fingerprinting...")
	cmd := exec.Command("ffmpeg", "-i", videoFile, "-vn", "-f", "f64le", "-acodec", "pcm_f64le", "-ac", "1", "-ar", "22050", "-")
	fingerprintOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to extract audio for fingerprinting: %v", err)
	}

	samples := make([]float64, len(fingerprintOutput)/8)
	for i := 0; i < len(samples); i++ {
		bits := uint64(0)
		for j := 0; j < 8; j++ {
			bits |= uint64(fingerprintOutput[i*8+j]) << (j * 8)
		}
		samples[i] = math.Float64frombits(bits)
	}

	fmt.Printf("✅ Extracted %d audio samples for fingerprinting\n", len(samples))

	fingerprint, err := GenerateFingerprint(samples)
	if err != nil {
		return nil, fmt.Errorf("failed to generate fingerprint: %v", err)
	}

	return fingerprint, nil
}

func startScreenShare(port string) *exec.Cmd {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ffmpeg", "-f", "gdigrab", "-i", "desktop", "-f", "dshow", "-i", "audio=Stereo Mix (2- Realtek(R) Audio)", "-vcodec", "libx264", "-acodec", "aac", "-f", "flv", fmt.Sprintf("rtmp://localhost:%s/live/stream", port))
	case "darwin":
		cmd = exec.Command("ffmpeg", "-f", "avfoundation", "-i", "1:0", "-vcodec", "libx264", "-acodec", "aac", "-f", "flv", fmt.Sprintf("rtmp://localhost:%s/live/stream", port))
	case "linux":
		cmd = exec.Command("ffmpeg", "-f", "x11grab", "-i", ":0.0", "-f", "pulse", "-i", "default", "-vcodec", "libx264", "-acodec", "aac", "-f", "flv", fmt.Sprintf("rtmp://localhost:%s/live/stream", port))
	}

	fmt.Println("Starting screen share stream...")
	cmd.Start()
	return cmd
}

func startHTTPServer(port string) *exec.Cmd {
	cmd := exec.Command("python", "-m", "http.server", port)
	fmt.Printf("Starting HTTP server on port %s\n", port)
	cmd.Start()
	return cmd
}

func ListAudioDevices() {
	cmd := exec.Command("ffmpeg", "-f", "dshow", "-list_devices", "true", "-i", "dummy")
	output, _ := cmd.CombinedOutput()
	fmt.Println("Available audio devices:")
	fmt.Println(string(output))
}
