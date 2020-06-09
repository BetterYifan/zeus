package main

import (
	"fmt"
	"zeus/model"
)

var (
	trie = *model.New("init trie")
)

func init() {
	trie.Insert("习近平", 10, "习近平")
	trie.Insert("温家宝", 10, "温家宝")
	trie.Insert("黄色", 10, "黄色")

	trie.Build()
}

func main() {
	hits := trie.Filter("我，习近平，打钱！")
	fmt.Println(len(hits))
}
