package rhymer

import (
	"reflect"
	"testing"

	"github.com/high-moctane/go-markov_chain_Japanese"
	"github.com/shogo82148/go-mecab"
)

func TestMorae(t *testing.T) {
	mecab, _ := mecab.New(map[string]string{})
	defer mecab.Destroy()
	parsed, _ := mecab.Parse("こんにちは世界")
	phrase := markov.MakePhraseString(parsed).Phrase()
	morae, _ := Morae(phrase)
	expected := []Mora{{"k", "o"}, {"*n", "*n"}, {"n", "i"}, {"ch", "i"}, {"w", "a"}, {"s", "e"}, {"k", "a"}, {"", "i"}}
	if !reflect.DeepEqual(expected, morae) {
		t.Errorf("expected %v, but %v", expected, morae)
	}
}

func TestSimilarity(t *testing.T) {
	var w MoraWeight
	var p0, p1 markov.Phrase
	mecab, _ := mecab.New(map[string]string{})
	var parsed string
	defer mecab.Destroy()
	m, _ := markov.New(1, map[string]string{})
	defer m.Destroy()

	parsed, _ = mecab.Parse("こんにちは")
	p0 = markov.MakePhraseString(parsed).Phrase()
	parsed, _ = mecab.Parse("こんにちは")
	p1 = markov.MakePhraseString(parsed).Phrase()
	w = NewMoraWeight([]MoraWeightCell{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}})
	if Similarity(p0, p1, w) != 1 {
		t.Errorf("expected %v, but %v", 1, Similarity(p0, p1, w))
	}

	parsed, _ = mecab.Parse("こんにちは")
	p0 = markov.MakePhraseString(parsed).Phrase()
	parsed, _ = mecab.Parse("魚市場")
	p1 = markov.MakePhraseString(parsed).Phrase()
	w = NewMoraWeight([]MoraWeightCell{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}})
	if Similarity(p0, p1, w) != 0.9393939393939394 {
		t.Errorf("expected %v, but %v", 0.9393939393939394, Similarity(p0, p1, w))
	}
}
