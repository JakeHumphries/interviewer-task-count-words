package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

// Part 1
// Questions:

// - What would be the time complexity of your solution?
// - very difficulty to answer because this is slow but here is an attempt:
// Splitting the text: If the text has ( n ) characters, the time complexity of strings.FieldsFunc is ( O(n) ).
// Iterating over words: If there are ( w ) words, iterating over them is ( O(w) ).
// Updating the frequency map: Updating the map for each word is ( O(1) ) on average, so for ( w ) words, it is ( O(w) ).
// Sorting the list: Sorting the list of word frequencies takes ( O(n log n) ), where ( n ) is the number of unique words.
// i dont know what this makes the overall time complexity....

// - What if we had to exclude specific words?
// implemented in the solution

// - What are the coding best practices you've used in your implementation?
// - seperation of concerns (each function is responsible for only one thing)
// - idiomatic error handling
// - case insensitivity
// - use of defer

var excludedWords = map[string]struct{}{
	"the":   {},
	"and":   {},
}

func main() {
	filePath := "example.txt"
	content, err := readFile(filePath)
	if err != nil {
		fmt.Println("failed to read file:", err)
		return
	}

	words := splitTextToWords(content)
	wf := countWordFrequencies(words, excludedWords)
	topWords := calcMostFrequentWords(wf, 10)

	fmt.Println("Top 10 most frequent words:")
	for _, wf := range topWords {
		fmt.Printf("%s: %d\n", wf.word, wf.count)
	}
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + " ")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}

func splitTextToWords(text string) []string {
	// splits the text via whitespace but we might want to use FieldsFunc if you need to handle punctuation and other non-alphanumeric characters
	// e.g. "Hello, world! This is a test. Testing, 1, 2, 3... Go is great; isn't it?"
	// output ["Hello," "world!" "This" "is" "a" "test." "Testing," "1," "2," "3..." "Go" "is" "great;" "isn't" "it?"]
	// example of fields func
	// 	words := strings.FieldsFunc(text, func(c rune) bool {
	//	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	// })

	// words := strings.Fields(text)

	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	return words
}

type wordFrequency struct {
	word  string
	count int
}

func countWordFrequencies(words []string, excludedWords map[string]struct{}) []wordFrequency {
	freqMap := make(map[string]int)
	for _, word := range words {
		// Convert to lowercase to handle potential casing issues
		word = strings.ToLower(word)
		// Skip the word if it is in the exclusion set
		if _, found := excludedWords[word]; found {
			continue
		}
		freqMap[word]++
	}

	var freqList []wordFrequency
	for word, count := range freqMap {
		freqList = append(freqList, wordFrequency{word, count})
	}

	return freqList
}

func calcMostFrequentWords(freqList []wordFrequency, amount int) []wordFrequency {
	if amount <= 0 {
		return []wordFrequency{}
	}

	// Sort the freqList slice in descending order based on the count field.
	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].count > freqList[j].count
	})

	// Select top words
	if amount > len(freqList) {
		amount = len(freqList)
	}
	return freqList[:amount]
}
