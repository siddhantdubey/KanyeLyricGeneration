package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

//lot of code from: https://rosettacode.org/wiki/Markov_chain_text_generator

type Chain struct {
	order       int                 //length of the prefix
	suffix      map[string][]string //suffixes of the prefix
	capitalized int                 //number of capitalized words in the chain
}

func isCapital(word string) bool {
	r, _ := utf8.DecodeRuneInString(word)
	return unicode.IsUpper(r)
}

func NewChain(reader io.Reader, order int) (*Chain, error) {
	chain := &Chain{
		order:  order,
		suffix: make(map[string][]string),
	}
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanWords)
	window := make([]string, order)
	for sc.Scan() {
		word := sc.Text()
		if len(window) > 0 {
			prefix := strings.Join(window, " ")
			chain.suffix[prefix] = append(chain.suffix[prefix], word)
			if isCapital(prefix) {
				chain.capitalized++
			}
		}
		window = append(window[1:], word)
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	return chain, nil
}

func NewChainFromFile(filename string, order int) (*Chain, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()
	return NewChain(file, order)
}

func (c *Chain) GenSentence(n int, startWithCap bool) string {
	var i int
	var words []string
	if startWithCap {
		i = rand.Intn(c.capitalized)
	} else {
		i = rand.Intn(len(c.suffix))
	}
	var prefix string
	for prefix = range c.suffix {
		if startWithCap && !isCapital(prefix) {
			continue
		}
		if i == 0 {
			break
		}
		i--
	}

	words = append(words, prefix)
	prefixWords := strings.Split(prefix, " ")
	n -= len(prefixWords)
	for {
		wordChoices := c.suffix[prefix]
		if len(wordChoices) == 0 {
			break
		}
		i = rand.Intn(len(wordChoices))
		suffix := wordChoices[i]
		words = append(words, suffix)

		n--
		if n < 0 || isSentenceEnd(suffix) {
			break
		}
		prefixWords = append(prefixWords[1:], suffix)
		prefix = strings.Join(prefixWords, " ")
	}
	return strings.Join(words, " ")
}

func isSentenceEnd(word string) bool {
	w, _ := utf8.DecodeLastRuneInString(word)
	return w == '.' || w == '?' || w == '!'
}

func main() {
	input := flag.String("in", "../data/lyrics/alllyrics.txt", "input file")
	n := flag.Int("n", 1, "number of words to use as prefix")
	runs := flag.Int("runs", 20, "number of runs to generate")
	wordsPerRun := flag.Int("words", 20, "number of words per run")
	startOnCapital := flag.Bool("capital", false, "start output with a capitalized prefix")
	flag.Parse()

	c, err := NewChainFromFile(*input, *n)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < *runs; i++ {
		out := c.GenSentence(*wordsPerRun, *startOnCapital)
		fmt.Println(out)
	}
}
