package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	t.Run("success case: it should read file content", func(t *testing.T) {
		// Create a temporary file with some content
		filePath := "testfile.txt"
		content := "Hello world! This is a test file."
		err := os.WriteFile(filePath, []byte(content), 0o644)

		assert.NoError(t, err)
		defer os.Remove(filePath)

		result, err := readFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, content+" ", result) // Note the extra space added by readFile
	})

	t.Run("failure case: it should return error for non-existent file", func(t *testing.T) {
		_, err := readFile("nonexistent.txt")
		assert.Error(t, err)
	})
}

func TestSplitTextToWords(t *testing.T) {
	t.Run("success case: it should split text into words (and not include punctuation)", func(t *testing.T) {
		text := "Hello, world! This is a test."
		expected := []string{"Hello", "world", "This", "is", "a", "test"}
		result := splitTextToWords(text)
		assert.Equal(t, expected, result)
	})
}

func TestCountWordFrequencies(t *testing.T) {
	t.Run("success case: it should count word frequencies (including capitals) and exclude specific words", func(t *testing.T) {
		words := []string{"hello", "world", "Hello", "test", "the", "and", "The"}
		expected := []wordFrequency{
			{"hello", 2},
			{"world", 1},
			{"test", 1},
		}

		var testExcludedWords = map[string]struct{}{
			"the":   {},
			"and":   {},
		}
		result := countWordFrequencies(words, testExcludedWords)
		assert.ElementsMatch(t, expected, result)
	})
}

func TestCalcMostFrequentWords(t *testing.T) {
	t.Run("success case: it should return top N most frequent words", func(t *testing.T) {
		freqList := []wordFrequency{
			{"hello", 2},
			{"world", 1},
			{"test", 1},
		}
		expected := []wordFrequency{
			{"hello", 2},
			{"world", 1},
		}
		result := calcMostFrequentWords(freqList, 2)
		assert.Equal(t, expected, result)
	})

	t.Run("edge case: it should handle amount greater than list length", func(t *testing.T) {
		freqList := []wordFrequency{
			{"hello", 2},
			{"world", 1},
		}
		expected := []wordFrequency{
			{"hello", 2},
			{"world", 1},
		}
		result := calcMostFrequentWords(freqList, 10)
		assert.Equal(t, expected, result)
	})

	t.Run("edge case: it should handle zero or negative amount", func(t *testing.T) {
		freqList := []wordFrequency{
			{"hello", 2},
			{"world", 1},
		}
		expected := []wordFrequency{}
		result := calcMostFrequentWords(freqList, 0)
		assert.Equal(t, expected, result)

		result = calcMostFrequentWords(freqList, -1)
		assert.Equal(t, expected, result)
	})
}
