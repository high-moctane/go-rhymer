package rhymer

import (
	"math/rand"
	"sync"

	"github.com/high-moctane/go-mecab_slice"
)

// Markov はマルコフ連鎖を行うための構造体です。
type Markov struct {
	Chain      []*Cell    // 小マルコフグラフつなげたものです。
	SliceLen   int        // Chain の List の最大長さです。
	CellSize   int        // 各 Cell に格納する文の最大個数です。
	MoraWeight MoraWeight // 韻を判定するときの重み付けです。
	Threshold  float64    // 韻を判定するときのしきい値です。
	rw         *sync.RWMutex
}

// Cell は Markov.Chain の Element となる小さなマルコフ連鎖のグラフです。
type Cell struct {
	Chain  map[mecabs.Morpheme]map[mecabs.Morpheme]int // map[prefix]map[suffix]count のように格納します。
	Morphs []*mecabs.Morpheme                          // 入力したすべての Morpheme のポインタを格納します。
	Count  int                                         // 入力した文の数です。
}

// NewCell は初期化された Cell のポインタを返します。
func NewCell(cellSize int) *Cell {
	return &Cell{
		Chain:  make(map[mecabs.Morpheme]map[mecabs.Morpheme]int),
		Morphs: make([]*mecabs.Morpheme, 0, cellSize),
	}
}

func (c *Cell) addMorph(prefix, suffix *mecabs.Morpheme) {
	m, ok := c.Chain[*prefix]
	if !ok {
		c.Chain[*prefix] = map[mecabs.Morpheme]int{*suffix: 1}
		return
	}
	_, ok = m[*suffix]
	if !ok {
		m[*suffix] = 1
		return
	}
	m[*suffix]++
}

// Add は Cell に Phrase をもとにした連鎖を追加します。
func (c *Cell) Add(phrase mecabs.Phrase) {
	prefix := mecabs.EOS

	for i := len(phrase); i >= 0; i-- {
		c.addMorph(&prefix, &phrase[i])
		c.Morphs = append(c.Morphs, &phrase[i])
		prefix = phrase[i]
	}
	c.addMorph(&prefix, &mecabs.BOS)
	c.Count++
}

//  NewMarkov は初期化された Markov を返します。
func NewMarkov(sliceLen, cellSize int, weight MoraWeight, thresh float64) Markov {
	return Markov{
		Chain:      []*Cell{NewCell(cellSize)},
		SliceLen:   sliceLen,
		CellSize:   cellSize,
		MoraWeight: weight,
		Threshold:  thresh,
		rw:         new(sync.RWMutex),
	}
}

// Add は phrase を Markov に追加します。
func (m *Markov) Add(phrase mecabs.Phrase) {
	if len(phrase) <= 0 {
		return
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	// list を更新しないといけない場合
	if m.Chain[len(m.Chain)-1].Count >= m.CellSize {
		// 最古の Cell を破棄しないといけない場合
		if len(m.Chain) >= m.SliceLen {
			m.Chain[0] = nil
			m.Chain = m.Chain[1:]
		}
		m.Chain = append(m.Chain, NewCell(m.CellSize))
	}

	m.Chain[len(m.Chain)-1].Add(phrase)
}

func (m *Markov) Morphs() []*mecabs.Morpheme {
	set := make(map[mecabs.Morpheme]*mecabs.Morpheme)
	for _, cell := range m.Chain {
		for _, morph := range cell.Morphs {
			set[*morph] = morph
		}
	}
	ans := make([]*mecabs.Morpheme, 0, len(set))
	for _, morph := range set {
		ans = append(ans, morph)
	}
	return ans
}

func (m *Markov) RandMorphs() []*mecabs.Morpheme {
	morphs := m.Morphs()
	for i := len(morphs) - 1; i >= 1; i-- {
		n := rand.Intn(i + 1)
		morphs[n], morphs[i] = morphs[i], morphs[n]
	}
	return morphs
}

func (m *Markov) RhymingMorphs(l int) []*mecabs.Morpheme {
	ans := make([]*mecabs.Morpheme, 0, l)
	for morphs := m.RandMorphs(); len(ans) < l || len(morphs) > 1; morphs = morphs[1:] {
		for i := 1; i < len(morphs); i++ {
			if m.MoraWeight.Similarity(mecabs.Phrase{*morphs[0]}, mecabs.Phrase{*morphs[1]}) < m.Threshold {
				continue
			}
			if len(ans) < 1 {
				ans = append(ans, morphs[0])
			}
			ans = append(ans, morphs[i])
			morphs[1], morphs[i] = morphs[i], morphs[1]
		}
	}
	return ans
}
