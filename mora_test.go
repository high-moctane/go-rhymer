package rhymer

import (
	"reflect"
	"testing"

	"github.com/high-moctane/go-mecab_slice"
)

func TestMorae(t *testing.T) {
	mecabs, err := mecabs.New(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	defer mecabs.Destroy()
	phrase, err := mecabs.ParseToPhrase("こんにちは")
	if err != nil {
		t.Fatal(err)
	}
	morae, ok := Morae(phrase)
	expected := []Mora{{"k", "o"}, {"*n", "*n"}, {"n", "i"}, {"ch", "i"}, {"w", "a"}}
	if !ok || !reflect.DeepEqual(expected, morae) {
		t.Errorf("expected %v, but %v", expected, morae)
	}
}

func TestSimilarity(t *testing.T) {
	var p0, p1 mecabs.Phrase
	w := NewMoraWeight([]MoraWeightCell{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}})
	mecabs, err := mecabs.New(map[string]string{})
	defer mecabs.Destroy()
	if err != nil {
		t.Fatal(err)
	}

	p0, err = mecabs.ParseToPhrase("こんにちは")
	if err != nil {
		t.Fatal(err)
	}
	p1, err = mecabs.ParseToPhrase("こんにちは")
	if err != nil {
		t.Fatal(err)
	}
	if w.Similarity(p0, p1) != 1 {
		t.Errorf("expected %v, but %v", 1, w.Similarity(p0, p1))
	}

	p0, err = mecabs.ParseToPhrase("こんにちは")
	if err != nil {
		t.Fatal(err)
	}
	p1, err = mecabs.ParseToPhrase("魚市場")
	if err != nil {
		t.Fatal(err)
	}
	if w.Similarity(p0, p1) != 0.9393939393939394 {
		t.Errorf("expected %v, but %v", 0.9393939393939394, w.Similarity(p0, p1))
	}
}
