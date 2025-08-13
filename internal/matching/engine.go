package matching

import (
	"fmt"
	"sort"
	"strconv"

	"Shazam/internal/audio"
	"Shazam/internal/database"
)

func FindBestMatch(db *database.DB, queryFingerprint *audio.AudioFingerprint) (*database.MatchResult, error) {
	fmt.Println("🔍 Searching database for best match...")

	songs, err := db.GetAllSongs()
	if err != nil {
		return nil, err
	}

	if len(songs) == 0 {
		return &database.MatchResult{IsMatch: false, Confidence: 0.0}, nil
	}

	fmt.Printf("📚 Comparing against %d songs in database...\n", len(songs))

	bestMatch := &database.MatchResult{
		IsMatch:    false,
		Confidence: 0.0,
		Song:       nil,
	}

	for i, song := range songs {
		refFingerprint := audio.ConvertSongToFingerprint(song)

		result := SlideHamming(refFingerprint, queryFingerprint)

		fmt.Printf("  [%d/%d] %s - %s: %.1f%%\n", i+1, len(songs),
			song.Artist, song.Title, result.Confidence*100)

		if result.Confidence > bestMatch.Confidence {
			bestMatch.Confidence = result.Confidence
			bestMatch.MatchOffset = result.MatchOffset
			bestMatch.IsMatch = result.IsMatch
			bestMatch.Song = song

			timePerSegment := 512.0 / 22050.0
			bestMatch.TimeInSong = float64(result.MatchOffset) * timePerSegment
		}
	}

	return bestMatch, nil
}

func GetTopMatches(db *database.DB, queryFingerprint *audio.AudioFingerprint, topN int) ([]*database.MatchResult, error) {
	songs, err := db.GetAllSongs()
	if err != nil {
		return nil, err
	}

	var results []*database.MatchResult

	for _, song := range songs {
		refFingerprint := &audio.AudioFingerprint{
			TrackID:      strconv.Itoa(song.ID),
			TrackName:    song.Title,
			Artist:       song.Artist,
			Fingerprint:  song.Fingerprint,
			HashSegments: song.HashSegments,
		}

		result := SlideHamming(refFingerprint, queryFingerprint)

		matchResult := &database.MatchResult{
			IsMatch:     result.IsMatch,
			Confidence:  result.Confidence,
			MatchOffset: result.MatchOffset,
			Song:        song,
			TimeInSong:  float64(result.MatchOffset) * (512.0 / 22050.0),
		}

		results = append(results, matchResult)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})

	if topN > len(results) {
		topN = len(results)
	}

	return results[:topN], nil
}
