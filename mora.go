package rhymer

import "github.com/high-moctane/go-mecab_slice"

// カタカナと Mora の対応を保存したものです。
var (
	katakana = map[string]Mora{
		"ア": Mora{"", "a"}, "イ": Mora{"", "i"}, "ウ": Mora{"", "u"}, "エ": Mora{"", "e"}, "オ": Mora{"", "o"},
		"カ": Mora{"k", "a"}, "キ": Mora{"k", "i"}, "ク": Mora{"k", "u"}, "ケ": Mora{"k", "e"}, "コ": Mora{"k", "o"},
		"サ": Mora{"s", "a"}, "シ": Mora{"sh", "i"}, "ス": Mora{"s", "u"}, "セ": Mora{"s", "e"}, "ソ": Mora{"s", "o"},
		"タ": Mora{"t", "a"}, "チ": Mora{"ch", "i"}, "ツ": Mora{"ts", "u"}, "テ": Mora{"t", "e"}, "ト": Mora{"t", "o"},
		"ナ": Mora{"n", "a"}, "ニ": Mora{"n", "i"}, "ヌ": Mora{"n", "u"}, "ネ": Mora{"n", "e"}, "ノ": Mora{"n", "o"},
		"ハ": Mora{"h", "a"}, "ヒ": Mora{"h", "i"}, "フ": Mora{"f", "u"}, "ヘ": Mora{"h", "e"}, "ホ": Mora{"h", "o"},
		"マ": Mora{"m", "a"}, "ミ": Mora{"m", "i"}, "ム": Mora{"m", "u"}, "メ": Mora{"m", "e"}, "モ": Mora{"m", "o"},
		"ヤ": Mora{"y", "a"}, "ユ": Mora{"y", "u"}, "ヨ": Mora{"y", "o"},
		"ラ": Mora{"r", "a"}, "リ": Mora{"r", "i"}, "ル": Mora{"r", "u"}, "レ": Mora{"r", "e"}, "ロ": Mora{"r", "o"},
		"ワ": Mora{"w", "a"}, "ヲ": Mora{"", "o"}, "ン": Mora{"*n", "*n"},
		"ガ": Mora{"g", "a"}, "ギ": Mora{"g", "i"}, "グ": Mora{"g", "u"}, "ゲ": Mora{"g", "e"}, "ゴ": Mora{"g", "o"},
		"ザ": Mora{"z", "a"}, "ジ": Mora{"j", "i"}, "ズ": Mora{"z", "u"}, "ゼ": Mora{"z", "e"}, "ゾ": Mora{"z", "o"},
		"ダ": Mora{"d", "a"}, "ヂ": Mora{"j", "i"}, "ヅ": Mora{"z", "u"}, "デ": Mora{"d", "e"}, "ド": Mora{"d", "o"},
		"バ": Mora{"b", "a"}, "ビ": Mora{"b", "i"}, "ブ": Mora{"b", "u"}, "ベ": Mora{"b", "e"}, "ボ": Mora{"b", "o"},
		"パ": Mora{"p", "a"}, "ピ": Mora{"p", "i"}, "プ": Mora{"p", "u"}, "ペ": Mora{"p", "e"}, "ポ": Mora{"p", "o"},
		"キャ": Mora{"ky", "a"}, "キュ": Mora{"ky", "u"}, "キョ": Mora{"ky", "o"},
		"シャ": Mora{"sh", "a"}, "シュ": Mora{"sh", "u"}, "ショ": Mora{"sh", "o"},
		"チャ": Mora{"ch", "a"}, "チュ": Mora{"ch", "u"}, "チョ": Mora{"ch", "o"},
		"ニャ": Mora{"ny", "a"}, "ニュ": Mora{"ny", "u"}, "ニョ": Mora{"ny", "o"},
		"ヒャ": Mora{"hy", "a"}, "ヒュ": Mora{"hy", "u"}, "ヒョ": Mora{"hy", "o"},
		"ミャ": Mora{"my", "a"}, "ミュ": Mora{"my", "u"}, "ミョ": Mora{"my", "o"},
		"リャ": Mora{"ry", "a"}, "リュ": Mora{"ry", "u"}, "リョ": Mora{"ry", "o"},
		"ギャ": Mora{"gy", "a"}, "ギュ": Mora{"gy", "u"}, "ギョ": Mora{"gy", "o"},
		"ジャ": Mora{"j", "a"}, "ジュ": Mora{"j", "u"}, "ジョ": Mora{"j", "o"},
		"ビャ": Mora{"by", "a"}, "ビュ": Mora{"by", "u"}, "ビョ": Mora{"by", "o"},
		"ピャ": Mora{"py", "a"}, "ピュ": Mora{"py", "u"}, "ピョ": Mora{"py", "o"},
		"ファ": Mora{"f", "a"}, "フィ": Mora{"f", "i"}, "フェ": Mora{"f", "e"}, "フォ": Mora{"f", "o"},
		"フュ": Mora{"fy", "u"},
		"ウィ": Mora{"w", "i"}, "ウェ": Mora{"w", "e"}, "ウォ": Mora{"w", "o"},
		"ヴァ": Mora{"v", "a"}, "ヴィ": Mora{"v", "i"}, "ヴェ": Mora{"v", "e"}, "ヴォ": Mora{"v", "o"},
		"ツァ": Mora{"ts", "a"}, "ツィ": Mora{"ts", "i"}, "ツェ": Mora{"ts", "e"}, "ツォ": Mora{"ts", "o"},
		"チェ": Mora{"ch", "e"}, "シェ": Mora{"sh", "e"}, "ジェ": Mora{"j", "e"},
		"ティ": Mora{"t", "i"}, "ディ": Mora{"d", "i"},
		"デュ": Mora{"d", "u"}, "トゥ": Mora{"t", "u"},
		"ッ": Mora{"*xtu", "*xtu"},
	}
)

// 発音情報を表す構造体です。
type Mora struct {
	Consonant string
	Vowel     string
}

// Mora の距離を比較するための重み付けをする構造体です。
type MoraWeightCell struct {
	Consonant float64
	Vowel     float64
}

// Mora 列の距離を比較するための重み付けをする構造体です。
type MoraWeight struct {
	Cells []MoraWeightCell
	Max   float64
}

// c をもとに MoraWeight を初期化して生成します。
func NewMoraWeight(c []MoraWeightCell) MoraWeight {
	var sum float64
	for _, v := range c {
		sum += v.Consonant
		sum += v.Vowel
	}
	return MoraWeight{Cells: c, Max: sum}
}

// mecabs.Phrase を Mora 列に変換します
// Phrase 内に発音が定義されていない Morpheme が存在する場合、
// bool が false になります。
func Morae(p mecabs.Phrase) ([]Mora, bool) {
	var pron string
	var ok bool
	if pron, ok = p.Pronounciation(); !ok {
		return []Mora{}, false
	}

	runes := []rune(pron + "*")
	ans := make([]Mora, 0, len(runes)-1)

	for i, end := 0, len(runes)-1; i < end; i++ {
		if mora, ok := katakana[string(runes[i:i+2])]; ok {
			ans = append(ans, mora)
			i++
		} else if mora, ok := katakana[string(runes[i])]; ok {
			ans = append(ans, mora)
		} else if runes[i] == 'ー' {
			newMora := Mora{"", ans[len(ans)-1].Vowel}
			ans = append(ans, newMora)
		} else {
			return []Mora{}, false
		}
	}
	return ans, true
}

// p0, p1 の発音の類似度を返します。
func (w *MoraWeight) Similarity(p0, p1 mecabs.Phrase) float64 {
	var morae0, morae1 []Mora
	var ok bool
	if morae0, ok = Morae(p0); !ok {
		return 0.0
	}
	if morae1, ok = Morae(p1); !ok {
		return 0.0
	}
	len0, len1, lenw := len(morae0), len(morae1), len(w.Cells)
	lenMin := min(len0, len1, lenw)

	var sum float64
	for i := 0; i < lenMin; i++ {
		if morae0[len0-1-i].Vowel == morae1[len1-1-i].Vowel {
			sum += w.Cells[lenw-1-i].Vowel
		}
		if morae0[len0-1-i].Consonant == morae1[len1-1-i].Consonant {
			sum += w.Cells[lenw-1-i].Consonant
		}
	}
	return sum / w.Max
}

func min(x, y, z int) int {
	if x > y {
		if y > z {
			return z
		} else {
			return y
		}
	} else {
		if x > z {
			return z
		} else {
			return x
		}
	}
}
