package matching

import (
	"math"

	"Shazam/internal/audio"
	"Shazam/internal/database"
)

func HammingHex(a, b string) int {
	if len(a) != len(b) {
		return math.MaxInt32
	}
	d := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d
}

func SlideHamming(reference, query *audio.AudioFingerprint) *database.MatchResult {
	ref := reference.HashSegments
	qry := query.HashSegments
	if len(qry) == 0 || len(ref) == 0 {
		return &database.MatchResult{IsMatch: false, Confidence: 0}
	}

	const maxNibbleMismatches = 8

	best := &database.MatchResult{IsMatch: false, Confidence: 0, MatchOffset: -1}
	maxOffset := len(ref) - len(qry)
	if maxOffset < 0 {
		maxOffset = 0
	}

	for offset := 0; offset <= maxOffset; offset++ {
		matches := 0
		for i := 0; i < len(qry) && offset+i < len(ref); i++ {
			if HammingHex(ref[offset+i], qry[i]) <= maxNibbleMismatches {
				matches++
			}
		}
		conf := float64(matches) / float64(len(qry))
		if conf > best.Confidence {
			best = &database.MatchResult{
				IsMatch:     conf >= 0.35,
				Confidence:  conf,
				MatchOffset: offset,
			}
		}
	}

	return best
}
