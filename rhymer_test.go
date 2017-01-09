package rhymer

import (
	"math/rand"
	"strings"
	"testing"

	markov "github.com/high-moctane/go-markov_chain_Japanese"
)

var sampleInput = strings.Split(`ある日の暮方の事である。一人の下人が、羅生門の下で雨やみを待っていた。広い門の下には、この男のほかに誰もいない。ただ、所々丹塗の剥げた、大きな円柱に、蟋蟀が一匹とまっている。羅生門が、朱雀大路にある以上は、この男のほかにも、雨やみをする市女笠や揉烏帽子が、もう二三人はありそうなものである。それが、この男のほかには誰もいない。何故かと云うと、この二三年、京都には、地震とか辻風とか火事とか饑饉とか云う災がつづいて起った。そこで洛中のさびれ方は一通りではない。旧記によると、仏像や仏具を打砕いて、その丹がついたり、金銀の箔がついたりした木を、路ばたにつみ重ねて、薪の料に売っていたと云う事である。洛中がその始末であるから、羅生門の修理などは、元より誰も捨てて顧る者がなかった。するとその荒れ果てたのをよい事にして、狐狸が棲すむ。盗人が棲む。とうとうしまいには、引取り手のない死人を、この門へ持って来て、棄てて行くと云う習慣さえ出来た。そこで、日の目が見えなくなると、誰でも気味を悪るがって、この門の近所へは足ぶみをしない事になってしまったのである。その代りまた鴉がどこからか、たくさん集って来た。昼間見ると、その鴉が何羽となく輪を描いて、高い鴟尾のまわりを啼きながら、飛びまわっている。ことに門の上の空が、夕焼けであかくなる時には、それが胡麻をまいたようにはっきり見えた。鴉は、勿論、門の上にある死人の肉を、啄みに来るのである。――もっとも今日は、刻限が遅いせいか、一羽も見えない。ただ、所々、崩れかかった、そうしてその崩れ目に長い草のはえた石段の上に、鴉の糞が、点々と白くこびりついているのが見える。下人は七段ある石段の一番上の段に、洗いざらした紺の襖の尻を据えて、右の頬に出来た、大きな面皰を気にしながら、ぼんやり、雨のふるのを眺めていた`, "。")

func init() {
	rand.Seed(0)
}

func TestGenerate(t *testing.T) {
	ma, _ := markov.New(1, map[string]string{})
	defer ma.Destroy()
	for _, v := range sampleInput {
		ma.Add(v)
	}
	weight := NewMoraWeight([]MoraWeightCell{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}})
	rh := New(ma, &weight, 0.8, 6)

	ph := rh.Generate(3)
	expected := []string{"鴉の下人は七段", "洛中の上にある石段", "一通りでは七段"}
	for i, v := range ph {
		if len(ph) != 3 || v.OriginalForm() != expected[i] {
			t.Errorf("expected %v, but %v", expected, ph)
		}
	}
}

func TestGenerateFromKana(t *testing.T) {
	ma, _ := markov.New(1, map[string]string{})
	defer ma.Destroy()
	for _, v := range sampleInput {
		ma.Add(v)
	}
	weight := NewMoraWeight([]MoraWeightCell{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}})
	rh := New(ma, &weight, 0.8, 6)

	ph := rh.GenerateFromKana(3, "アンガイ")
	expected := []string{"その荒れ果てたと云う災", "ある日の暮方のまわり", "するとその丹がその代り"}
	for i, v := range ph {
		if len(ph) != 3 || v.OriginalForm() != expected[i] {
			t.Errorf("expected %v, but %v", expected, ph)
		}
	}
}

func TestStream(t *testing.T) {
	ma, _ := markov.New(1, map[string]string{})
	defer ma.Destroy()
	for _, v := range sampleInput {
		ma.Add(v)
	}
	weight := NewMoraWeight([]MoraWeightCell{{1.0, 10.0}, {2.0, 20.0}, {3.0, 30.0}})
	rh := New(ma, &weight, 0.8, 6)

	ch, kill := rh.Stream(3)
	defer func() { kill <- true }()

	expected := [][]string{
		{"下人は七段ある石段", "盗人が一人の段", "するとその鴉は七段"},
		{"下人は七段ある石段", "そこで洛中の襖の一番", "羅生門の下にある石段"},
	}

	for i := 0; i < 2; i++ {
		for j, v := range <-ch {
			if v.OriginalForm() != expected[i][j] {
				t.Errorf("expected %v, but %v", expected[i][j], v.OriginalForm())
			}
		}
	}
}
