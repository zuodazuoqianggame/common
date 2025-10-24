package utils

import (
	"bufio"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var tire = NewTrie()

func init() {
	f, err := os.Open("./config/trie.txt")
	defer f.Close()
	if err != nil {
		log.Error(err)
		return
	}
	buff := bufio.NewReader(f)
	for {
		line, _, err := buff.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
		}
		tire.Inster(string(line))
	}

}

func TrieReplace(str string) string {
	return tire.Replace(str)
}
