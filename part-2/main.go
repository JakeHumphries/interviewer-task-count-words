package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"
)

// Part 2

// use chunking along with go routines and channels to achieve concurrency
// Questions:

// - If the candidate is using a Go map ask if it's safe for concurrency
// - Go maps are not safe for concurrent read and write operations. If multiple goroutines access a map concurrently without proper synchronization, it can lead to race conditions and undefined behavior.
// - To make the map safe for concurrent access, you can use synchronization mechanisms provided by Go, such as sync.Mutex or sync.RWMutex. Alternatively, you can use the sync.Map type, which is designed for concurrent use.

// - What if we wanted to limit the amount of CPU threads used at a time (bounded concurrency)?
// - in GO we can use a buffered channel to control the number of concurrent goroutines. (make(chan struct{}, maxGoroutines)

// - What are some potential race conditions with this implementation?
// the main race condition is fixed by using mutexs.
// The current implementation locks the entire finalFreqMap for each chunk. This can lead to contention and reduced performance if many goroutines are trying to update the map simultaneously

// bonus question: ask the person about error handling with goroutines:
// The current implementation does not handle errors within goroutines. If a goroutine encounters an error, it will silently fail, potentially leading to incomplete or incorrect results.

var excludedWords = map[string]struct{}{
	"the": {},
	"and": {},
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
	// Split the text into chunks
	chunkSize := 1000 // Adjust this size
	var chunks []string

	// 	Initialization (i := 0): The loop starts with i initialized to 0.
	// Condition (i < len(text)): The loop continues to run as long as i is less than the length of the text.
	// Increment (i += chunkSize): After each iteration, i is incremented by chunkSize.
	// 	Calculate end (end := i + chunkSize): The end variable is set to the current value of i plus chunkSize. This determines the end index of the current chunk.
	// Adjust end if it exceeds the text length:
	// If end is greater than the length of the text (end > len(text)), it is adjusted to be the length of the text (end = len(text)). This ensures that the last chunk does not go out of bounds.
 	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	// Channel to collect words from each chunk
	wordChan := make(chan []string, len(chunks))
	var wg sync.WaitGroup

	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk string) {
			defer wg.Done()
			words := strings.FieldsFunc(chunk, func(c rune) bool {
				return !unicode.IsLetter(c) && !unicode.IsNumber(c)
			})
			wordChan <- words
		}(chunk)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(wordChan)
	}()

	// Collect all words from the channel
	var allWords []string
	for words := range wordChan {
		allWords = append(allWords, words...)
	}

	return allWords
}

type wordFrequency struct {
	word  string
	count int
}

func countWordFrequencies(words []string, excludedWords map[string]struct{}) []wordFrequency {
	// Split the words into chunks
	chunkSize := 1000 // Adjust this size based on your needs
	var chunks [][]string
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, words[i:end])
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	finalFreqMap := make(map[string]int)

	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			freqMap := make(map[string]int)
			for _, word := range chunk {
				// Convert to lowercase to handle potential casing issues
				word = strings.ToLower(word)
				// Skip the word if it is in the exclusion set
				if _, found := excludedWords[word]; found {
					continue
				}
				freqMap[word]++
			}
			mu.Lock()
			for word, count := range freqMap {
				finalFreqMap[word] += count
			}
			mu.Unlock()
		}(chunk)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	var freqList []wordFrequency
	for word, count := range finalFreqMap {
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

// SOLUTION USING SYNC MAP INSTEAD OF MUTEX

// func countWordFrequenciesSynMap(words []string, excludedWords map[string]struct{}) []wordFrequency {
// 	// Split the words into chunks
// 	chunkSize := 1000 // Adjust this size based on your needs
// 	var chunks [][]string
// 	for i := 0; i < len(words); i += chunkSize {
// 		end := i + chunkSize
// 		if end > len(words) {
// 			end = len(words)
// 		}
// 		chunks = append(chunks, words[i:end])
// 	}

// 	// Channel to collect frequency maps from each chunk
// 	var wg sync.WaitGroup
// 	finalFreqMap := sync.Map{}

// 	for _, chunk := range chunks {
// 		wg.Add(1)
// 		go func(chunk []string) {
// 			defer wg.Done()
// 			freqMap := make(map[string]int)
// 			for _, word := range chunk {
// 				// Convert to lowercase to handle potential casing issues
// 				word = strings.ToLower(word)
// 				// Skip the word if it is in the exclusion set
// 				if _, found := excludedWords[word]; found {
// 					continue
// 				}
// 				freqMap[word]++
// 			}
// 			for word, count := range freqMap {
// 				actual, _ := finalFreqMap.LoadOrStore(word, count)
// 				if actual != count {
// 					finalFreqMap.Store(word, actual.(int)+count)
// 				}
// 			}
// 		}(chunk)
// 	}

// 	// Wait for all goroutines to finish
// 	wg.Wait()

// 	var freqList []wordFrequency
// 	finalFreqMap.Range(func(key, value interface{}) bool {
// 		freqList = append(freqList, wordFrequency{key.(string), value.(int)})
// 		return true
// 	})

// 	return freqList
// }


