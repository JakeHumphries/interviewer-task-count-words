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
// n = number of lines in the file, m = length of each line in the file, k = total number of characters in the entire text, w = total number of words in the text
// ReadFile: Opening the file: ( O(1) )
// Reading the file: The file is read line by line. If there are ( n ) lines and each line has ( m ) characters on average, the time complexity is ( O(n times m) ).
// Appending to strings.Builder: Appending each line to the strings.Builder is also ( O(n times m) ).
// Splitting the text: If the text has ( k ) characters, the time complexity of strings.FieldsFunc is ( O(k) ).
// Iterating over words: If there are ( w ) words, iterating over them is ( O(w) ).
// Updating the frequency map: Updating the map for each word is ( O(1) ) on average, so for ( w ) words, it is ( O(w) ).
// Sorting the list: Sorting the list of word frequencies takes ( O(w log w) ), where ( w ) is the number of unique words.
// Selecting top words: Selecting the top amount words is ( O(text{amount}) ).
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
