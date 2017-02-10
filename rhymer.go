package rhymer

import (
	"math/rand"
	"time"

	"github.com/high-moctane/go-mecab_slice"
)

type Rhymer struct {
	Markov  Markov
	Length  int
	Lines   int
	Weight  MoraWeight
	Thresh  float64
	Timeout time.Duration
}

func New(
	markov Markov,
	length, lines int,
	weight MoraWeight,
	thresh float64,
	to time.Duration,
) Rhymer {
	return Rhymer{
		Markov:  markov,
		Length:  length,
		Lines:   lines,
		Weight:  weight,
		Thresh:  thresh,
		Timeout: to,
	}
}

func shuffle(morphs []*mecabs.Morpheme) []*mecabs.Morpheme {
	for i := len(morphs) - 1; i >= 1; i-- {
		n := rand.Intn(i)
		morphs[i], morphs[n] = morphs[n], morphs[i]
	}
	return morphs
}

func (r *Rhymer) RhymingMorphs() []*mecabs.Morpheme {
}

func (r *Rhymer) isRhyme(ph0, ph1 mecabs.Phrase) bool {
	if r.Weight.Similarity(ph0, ph1) < r.Thresh {
		return false
	}
	return true
}
