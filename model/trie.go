package model

import (
	"container/list"
	"fmt"
	"strings"
)

type trieNode struct {
	isEnd bool
	index int
	fail  *trieNode
	child map[rune]*trieNode

	//extra data
	level int8
	data  interface{}
}

func newNode() *trieNode {
	return &trieNode{
		fail:  nil,
		index: -1,
		child: make(map[rune]*trieNode),
	}
}

// New .
func New(desc string) *Trie {
	return &Trie{
		Size:         0,
		root:         newNode(),
		indexSizeMap: make(map[int]int),
		Words:        make([]string, 0),
		Desc:         desc,
	}
}

// Trie ac search trie
type Trie struct {
	Size         int //包含多少词条
	root         *trieNode
	Words        []string
	Desc         string
	indexSizeMap map[int]int // map[index]size 对应index的词条byte长度
}

func (t Trie) String() string {
	if len(t.Words) <= 20 {
		return fmt.Sprintf("Trie: %s, words: %+v", t.Desc, t.Words)
	}
	return fmt.Sprintf("Trie: %s, words: %+v...", t.Desc, t.Words[:20])
}

// Insert insert a node
func (t *Trie) Insert(s string, level int8, data interface{}) {
	curNode := t.root
	for _, v := range s {
		if curNode.child[v] == nil {
			curNode.child[v] = newNode()
		}
		curNode = curNode.child[v]
	}
	curNode.index = t.Size
	curNode.isEnd = true

	// extra data
	t.indexSizeMap[t.Size] = len([]byte(s))
	curNode.level = level
	curNode.data = data
	t.Size++
	t.Words = append(t.Words, s)
}

// Build build trie
func (t *Trie) Build() {
	ll := list.New()
	ll.PushBack(t.root)
	for ll.Len() > 0 {
		curNode := ll.Remove(ll.Front()).(*trieNode)
		var p *trieNode
		for r, childNode := range curNode.child {
			if curNode == t.root {
				childNode.fail = t.root
			} else {
				p = curNode.fail
				for p != nil {
					if p.child[r] != nil {
						childNode.fail = p.child[r]
						break
					}
					p = p.fail
				}
				if p == nil {
					childNode.fail = t.root
				}
			}
			ll.PushBack(childNode)
		}
	}
}

// Hit .
type Hit struct {
	Level int8
	Pos   []int
	Extra interface{}
}

// Filter .
func (t *Trie) Filter(in string) (hits []*Hit) {
	var (
		curNode  = t.root
		p        *trieNode
		byteSize = 0
	)
	for _, r := range strings.ToLower(in) {
		byteSize += len([]byte(string(r)))
		for curNode.child[r] == nil && curNode != t.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[r]
		if curNode == nil {
			curNode = t.root
		}
		p = curNode
		if p != t.root && p.isEnd {
			hits = append(hits, &Hit{
				Level: p.level,
				Pos:   []int{byteSize - t.indexSizeMap[p.index], byteSize},
				Extra: p.data,
			})
			p = p.fail
		}
	}
	return
}
