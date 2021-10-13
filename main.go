package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

func main() {
	defer elapsedTime()()
	strTarget := "a-"

	// all even values start with 'a-' and all odd values start with 'b-'
	wordList := genValues(50)
	firstPartSize, lastPartSize := splitSizes(len(wordList))

	if len(wordList) <= 100 {
		fmt.Printf("Word List: %v \n\n", wordList)
	}

	// fetch all words start with srtTarget
	retList := findStartWith(wordList, strTarget, firstPartSize, lastPartSize)
	sort.Strings(retList)

	fmt.Printf("Sugestion words: %v \n", retList)
}

func findStartWith(wordList []string, target string, size1, size2 int) []string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			findFromBegin(wordList, target, size1, ch)
			wg.Done()
		}()

		go func() {
			findFromEnd(wordList, target, size2, ch)
			wg.Done()
		}()
		wg.Wait()
	}()

	out := make(chan []string)
	go func(out chan []string, ch1 chan string) {
		defer close(out)
		var retList []string
		for v := range ch1 {
			retList = append(retList, v)
		}
		out <- retList
	}(out, ch)

	return <-out
}

func findFromBegin(wordList []string, target string, size int, ch1 chan string) {
	for i := 0; i <= size; i++ {
		word := wordList[i]
		if len(word) < len(target) {
			continue
		}

		if matchWord(word, target) {
			ch1 <- word
		}
	}
}

func findFromEnd(wordList []string, target string, size int, ch chan string) {
	for i := len(wordList) - 1; i >= size; i-- {
		word := wordList[i]
		if len(word) < len(target) {
			continue
		}

		if matchWord(word, target) {
			ch <- word
		}
	}
}

func matchWord(word, target string) bool {
	pw := word[0:len(target)]
	return strings.EqualFold(pw, target)
}

func splitSizes(listSize int) (firstPartSize, lastPartSize int) {
	if listSize == 0 {
		return 0, 0
	}

	size := 0
	if listSize%2 == 0 {
		size = listSize / 2
		return size, size
	}

	// return values to odd slices
	size = (listSize - 1) / 2
	return size, size + 1
}

func genValues(size int) []string {
	fmt.Println("Starting genValues....")
	wordList := make([]string, size)
	for i := 0; i < size; i++ {
		w := ""
		if i%2 == 0 {
			w = fmt.Sprintf("a-%d", i)
		} else {
			w = fmt.Sprintf("b-%d", i)
		}
		wordList[i] = w
	}
	fmt.Println("Finished genValues....")
	return wordList
}

func elapsedTime() func() {
	start := time.Now()
	return func() {
		fmt.Printf("Elapsed Time: %v", time.Since(start))
	}
}
