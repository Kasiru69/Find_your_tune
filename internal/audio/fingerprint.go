package audio

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/cmplx"
	"sort"
	"strings"

	"Shazam/internal/database"
)

type AudioFingerprint struct {
	TrackID      string   `json:"track_id"`
	TrackName    string   `json:"track_name"`
	Artist       string   `json:"artist"`
	Fingerprint  string   `json:"fingerprint"`
	HashSegments []string `json:"hash_segments"`
}

func GenerateFingerprint(samples []float64) (*AudioFingerprint, error) {
	if len(samples) < 1024 {
		return nil, fmt.Errorf("insufficient audio samples")
	}

	windowSize := 2048
	hopSize := 512
	hashSegments := []string{}

	for i := 0; i < len(samples)-windowSize; i += hopSize {
		window := samples[i : i+windowSize]

		if rms(window) < 0.00001 {
			continue
		}

		windowed := make([]complex128, windowSize)
		for j, sample := range window {
			w := 0.54 - 0.46*math.Cos(2*math.Pi*float64(j)/float64(windowSize-1))
			windowed[j] = complex(sample*w, 0)
		}

		spectrum := fft(windowed)
		hash := createRobustHash(spectrum)
		if hash != "" {
			hashSegments = append(hashSegments, hash)
		}
	}

	if len(hashSegments) < 1 {
		return nil, fmt.Errorf("too few hash segments generated: %d", len(hashSegments))
	}

	combinedHash := strings.Join(hashSegments, "")
	finalHash := sha256.Sum256([]byte(combinedHash))

	return &AudioFingerprint{
		Fingerprint:  hex.EncodeToString(finalHash[:]),
		HashSegments: hashSegments,
	}, nil
}

func createRobustHash(spectrum []complex128) string {
	n := len(spectrum) / 2
	mags := make([]float64, n)
	const eps = 1e-12
	for i := 0; i < n; i++ {
		mags[i] = cmplx.Abs(spectrum[i]) + eps
	}

	numBands := 16
	bandSize := n / numBands
	bands := make([]float64, numBands)

	for b := 0; b < numBands; b++ {
		start := b * bandSize
		end := start + bandSize
		if b == numBands-1 || end > n {
			end = n
		}
		var sum float64
		for i := start; i < end; i++ {
			sum += mags[i]
		}
		avg := sum / float64(end-start)

		bands[b] = 20.0 * math.Log10(avg)
	}

	m := median(bands)
	for b := range bands {
		bands[b] -= m
	}

	const lo, hi = -24.0, 24.0
	hash := make([]byte, numBands)
	hexDigits := []byte("0123456789abcdef")
	for b := 0; b < numBands; b++ {
		v := bands[b]
		if v < lo {
			v = lo
		}
		if v > hi {
			v = hi
		}
		q := int(math.Round((v - lo) / (hi - lo) * 15.0))
		hash[b] = hexDigits[q]
	}
	return string(hash)
}

func median(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	tmp := append([]float64(nil), xs...)
	sort.Float64s(tmp)
	n := len(tmp)
	if n%2 == 1 {
		return tmp[n/2]
	}
	return 0.5 * (tmp[n/2-1] + tmp[n/2])
}

func fft(x []complex128) []complex128 {
	n := len(x)
	if n <= 1 {
		return x
	}
	even := make([]complex128, n/2)
	odd := make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		even[i] = x[2*i]
		odd[i] = x[2*i+1]
	}
	evenFFT := fft(even)
	oddFFT := fft(odd)
	result := make([]complex128, n)
	for i := 0; i < n/2; i++ {
		t := cmplx.Exp(complex(0, -2*math.Pi*float64(i)/float64(n))) * oddFFT[i]
		result[i] = evenFFT[i] + t
		result[i+n/2] = evenFFT[i] - t
	}
	return result
}

func rms(samples []float64) float64 {
	var sum float64
	for _, sample := range samples {
		sum += sample * sample
	}
	return math.Sqrt(sum / float64(len(samples)))
}

func ConvertSongToFingerprint(song *database.Song) *AudioFingerprint {
	return &AudioFingerprint{
		TrackID:      fmt.Sprintf("%d", song.ID),
		TrackName:    song.Title,
		Artist:       song.Artist,
		Fingerprint:  song.Fingerprint,
		HashSegments: song.HashSegments,
	}
}
