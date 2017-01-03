package rhymer

import (
	"reflect"
	"runtime"

	markov "github.com/high-moctane/go-markov_chain_Japanese"
)

type Rhymer struct {
	Markov     *markov.Markov
	Weight     *MoraWeight
	MaxLen     int
	MaxPhrase  int
	Similarity float64
}

func New(m *markov.Markov, w *MoraWeight, max int, maxp int, sim float64) Rhymer {
	return Rhymer{Markov: m, Weight: w, MaxLen: max, MaxPhrase: maxp, Similarity: sim}
}

func endCondition(m markov.Morpheme) bool {
	return false
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

	// if reflect.DeepEqual(p0[len(p0)-1], p1[len(p1)-1]) {
	// return false
	// }

	return true
}

func isDup(ph []markov.Phrase) bool {
	for i := 0; i < len(ph)-1; i++ {
		for j := i + 1; j < len(ph); j++ {
			if reflect.DeepEqual(ph[i][len(ph[i])-1], ph[j][len(ph[j])-1]) {
				return true
			}
		}
	}
	return false
}

func (r *Rhymer) TryGenerate() ([]markov.Phrase, bool) {
	ph := make([]markov.Phrase, r.MaxPhrase)
	ans := make([]markov.Phrase, 0, r.MaxPhrase)
	for i, _ := range ph {
		ph[i] = r.Markov.Generate(r.MaxLen, endCondition)
	}
	ans = append(ans, ph[0])
	for i := 1; i < len(ph); i++ {
		if r.isRhyme(ph[0], ph[i]) {
			ans = append(ans, ph[i])
		}
	}
	if len(ans) <= 1 {
		return []markov.Phrase{}, false
	}
	if isDup(ans) {
		return []markov.Phrase{}, false
	}
	return ans, true
}

func (r *Rhymer) Generate() []markov.Phrase {
	for {
		if ph, ok := r.TryGenerate(); ok {
			return ph
		}
		runtime.Gosched()
	}
}

func (r *Rhymer) Stream() <-chan []markov.Phrase {
	ch := make(chan []markov.Phrase, 1+runtime.NumCPU())
	for i := 0; i < runtime.NumCPU()+1; i++ {
		go func() {
			for {
				ch <- r.Generate()
				runtime.Gosched()
			}
		}()
	}
	return (<-chan []markov.Phrase)(ch)
}
