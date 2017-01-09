package rhymer

import (
	"runtime"

	markov "github.com/high-moctane/go-markov_chain_Japanese"
)

type Rhymer struct {
	Markov     *markov.Markov
	Weight     *MoraWeight
	Similarity float64
	Morphemes  int
}

func endCondition(_ markov.Morpheme) bool {
	return false
}

func New(m *markov.Markov, w *MoraWeight, s float64, mo int) Rhymer {
	return Rhymer{Markov: m, Weight: w, Similarity: s, Morphemes: mo}
}

func (r *Rhymer) isRhyme(p0, p1 markov.Phrase) bool {
	if Similarity(p0, p1, *r.Weight) < r.Similarity {
		return false
	}

	for _, v := range []string{"助詞", "形容詞", "動詞", "助動詞", "記号", "感動詞", "連体詞", "副詞", "接頭詞"} {
		if p0[len(p0)-1].PartOfSpeech == v {
			return false
		}
		if p1[len(p1)-1].PartOfSpeech == v {
			return false
		}
	}

	for _, v := range []string{"非自立"} {
		if p0[len(p0)-1].PartOfSpeechSection1 == v {
			return false
		}
		if p1[len(p1)-1].PartOfSpeechSection1 == v {
			return false
		}
	}

	return true
}

func (r *Rhymer) isDup(ph []markov.Phrase) bool {
	for i := 0; i < len(ph)-1; i++ {
		for j := i + 1; j < len(ph); j++ {
			if ph[i][len(ph[i])-1].OriginalForm == ph[j][len(ph[j])-1].OriginalForm {
				return true
			}
		}
	}
	return false
}

func (r *Rhymer) GeneratePair() []markov.Phrase {
	buf := make([]markov.Phrase, 2)
	shift := func() {
		buf[1] = buf[0]
		buf[0] = r.Markov.Generate(r.Morphemes, endCondition)
	}
	shift()

	for {
		shift()
		if r.isRhyme(buf[0], buf[1]) && !r.isDup(buf) {
			return buf
		}
	}
}

func (r *Rhymer) GenerateFromPhrase(l int, p markov.Phrase) []markov.Phrase {
	for {
		ans := make([]markov.Phrase, 1, l+1)
		ans[0] = p
		for len(ans) < l+1 {
			for {
				ph := r.Markov.Generate(r.Morphemes, endCondition)
				if r.isRhyme(ph, ans[len(ans)-1]) {
					ans = append(ans, ph)
					break
				}
			}
		}

		if !r.isDup(ans) {
			return ans[1:]
		}
	}
}

func (r *Rhymer) GenerateFromKana(l int, s string) []markov.Phrase {
	p := markov.Phrase{{Pronounciation: s}}
	return r.GenerateFromPhrase(l, p)
}

func (r *Rhymer) Generate(l int) []markov.Phrase {
	pair := r.GeneratePair()
	return append(pair, r.GenerateFromPhrase(l-2, pair[1])...)
}

func (r *Rhymer) Stream(l int) (<-chan []markov.Phrase, chan<- bool) {
	kill := make(chan bool, 1)
	ans := make(chan []markov.Phrase, 1)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-kill:
					return
				default:
					ans <- r.Generate(l)
				}
			}
		}()
	}
	return (<-chan []markov.Phrase)(ans), (chan<- bool)(kill)
}
