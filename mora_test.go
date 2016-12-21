package rhymer

import (
	"reflect"
	"testing"
)

func TestParseToMora(t *testing.T) {
	var expected [][]string

	if mora, ok := parseToMora("ポワ"); ok {
		expected = [][]string{{"p", "o"}, {"w", "a"}}
		if !reflect.DeepEqual(mora, expected) {
			t.Errorf("expected %v, but %v", expected, mora)
		}
	} else {
		t.Errorf("parse error.")
	}

	if mora, ok := parseToMora("チョコ"); ok {
		expected = [][]string{{"ch", "o"}, {"k", "o"}}
		if !reflect.DeepEqual(mora, expected) {
			t.Errorf("expected %v, but %v", expected, mora)
		}
	} else {
		t.Errorf("parse error.")
	}

	if mora, ok := parseToMora("ポワー"); ok {
		expected = [][]string{{"p", "o"}, {"w", "a"}, {"", "a"}}
		if !reflect.DeepEqual(mora, expected) {
			t.Errorf("expected %v, but %v", expected, mora)
		}
	} else {
		t.Errorf("parse error.")
	}

	if mora, ok := parseToMora("ジュース"); ok {
		expected = [][]string{{"j", "u"}, {"", "u"}, {"s", "u"}}
		if !reflect.DeepEqual(mora, expected) {
			t.Errorf("expected %v, but %v", expected, mora)
		}
	} else {
		t.Errorf("parse error.")
	}

	if _, ok := parseToMora("ポワpowa"); ok {
		t.Errorf("parse error.")
	}
}
