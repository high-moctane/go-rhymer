package rhymer

import (
	"sync"

	"github.com/high-moctane/go-mecab_slice"
)

// Markov はマルコフ連鎖を行うための構造体です。
// NOTE:
//  Chain はインデックスが若いほうが古い
type Markov struct {
	Chain     []*Cell // 小マルコフグラフつなげたものです。
	CellCount int     // Chain の List の最大長さです。
	CellSize  int     // 各 Cell に格納する文の最大個数です。
	rw        *sync.RWMutex
}

// Cell は Markov.Chain の Element となる小さなマルコフ連鎖のグラフです。
type Cell struct {
	Chain  map[*mecabs.Morpheme]map[*mecabs.Morpheme]int // map[suffix]map[prefix]count のように格納します。
	Morphs map[mecabs.Morpheme]*mecabs.Morpheme          // 入力したすべての Morpheme のポインタを格納します。
	Count  int                                           // 入力した文の数です。
}

// NewCell は初期化された Cell のポインタを返します。
func NewCell(cellSize int) *Cell {
	return &Cell{
		Chain:  make(map[*mecabs.Morpheme]map[*mecabs.Morpheme]int),
		Morphs: make(map[mecabs.Morpheme]*mecabs.Morpheme),
	}
}

// Pointer は morph のポインタを返します。
func (c *Cell) Pointer(morph *mecabs.Morpheme) *mecabs.Morpheme {
	p, ok := c.Morphs[*morph]
	if !ok {
		c.Morphs[*morph] = morph
		return morph
	}
	return p
}

// AddPair は initial から terminal へ伸びる辺を追加します。
func (c *Cell) AddPair(initial, terminal *mecabs.Morpheme) {
	pInit := c.Pointer(initial)
	pTerm := c.Pointer(terminal)
	ma, ok := c.Chain[pInit]
	if !ok {
		c.Chain[pInit] = map[*mecabs.Morpheme]int{pTerm: 1}
		return
	}
	_, ok = ma[pTerm]
	if !ok {
		ma[pTerm] = 1
	}
	ma[pTerm]++
}

// Add は cell に phrase を登録します。
func (c *Cell) Add(phrase mecabs.Phrase) {
	initial := &mecabs.EOS

	for i := len(phrase) - 1; i >= 0; i-- {
		terminal := &phrase[i]
		c.AddPair(initial, terminal)
		initial = terminal
	}
	c.AddPair(initial, &mecabs.BOS)
	c.Count++
}

// NewMarkov は初期化された Markov を返します。
func NewMarkov(cellCount, cellSize int) Markov {
	return Markov{
		Chain:     []*Cell{NewCell(cellSize)},
		CellCount: cellCount,
		CellSize:  cellSize,
		rw:        new(sync.RWMutex),
	}
}

// Add は Markov に Phrase を追加します。
func (m *Markov) Add(phrase mecabs.Phrase) {
	m.rw.Lock()
	defer m.rw.Unlock()

	// Cell の更新が必要な場合
	if m.Chain[len(m.Chain)-1].Count >= m.CellSize {
		// 古い Cell を削除しないといけない場合
		if len(m.Chain) >= m.CellCount {
			m.Chain[0] = nil // GC に食わせるため
			m.Chain = m.Chain[1:]
		}
		m.Chain = append(m.Chain, NewCell(m.CellSize))
	}
	m.Chain[len(m.Chain)-1].Add(phrase)
}

// Morphs は Markov に格納されたすべての Morpheme のポインタを返します。
// 重複を許しません。
func (m *Markov) Morphs() []*mecabs.Morpheme {
	set := make(map[mecabs.Morpheme]*mecabs.Morpheme)
	ans := make([]*mecabs.Morpheme, 0, len(set))

	m.rw.RLock()
	defer m.rw.RUnlock()

	for _, cell := range m.Chain {
		for morph, p := range cell.Morphs {
			set[morph] = p
		}
	}
	for _, p := range set {
		ans = append(ans, p)
	}
	return ans
}

// Next は morph の次にくる可能性のある Morpheme のスライスを返します。
func (m *Markov) Next(morph *mecabs.Morpheme) []*mecabs.Morpheme {
	set := make(map[mecabs.Morpheme]*mecabs.Morpheme)
	ans := make([]*mecabs.Morpheme, 0, len(set))

	m.rw.RLock()
	defer m.rw.RUnlock()

	for _, cell := range m.Chain {
		candidates, ok := cell.Chain[cell.Pointer(morph)]
		if !ok {
			continue
		}
		for mo := range candidates {
			set[*mo] = mo
		}
	}
	for _, p := range set {
		ans = append(ans, p)
	}
	return ans
}
