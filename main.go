package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	sastrawi "github.com/RadhiFadlillah/go-sastrawi"
	"github.com/jbrukh/bayesian"
)

var (
	stopwords  sastrawi.Dictionary
	dictionary sastrawi.Dictionary
	queryWords sastrawi.Dictionary
	stemmer    sastrawi.Stemmer

	negativeWords []string
	positiveWords []string
	s             []string
	classifier    *bayesian.Classifier
)

const (
	Good bayesian.Class = "Good"
	Bad  bayesian.Class = "Bad"
)

var (
	dictionaryQueryWords = []string{"goto", "ipo", "gojek", "tokopedia", "tokped", "toped"}
)

func init() {

	stopwords = sastrawi.DefaultStopword()
	dictionary = sastrawi.DefaultDictionary()
	stemmer = sastrawi.NewStemmer(dictionary)
	queryWords = sastrawi.NewDictionary(dictionaryQueryWords...)

	// if 0 then Bad
	// if 1 then Good
	classifier = bayesian.NewClassifier(Bad, Good)

}

func main() {
	t := time.Now()

	// 1. Read
	s = ReadSourceReadonly("twit_translated_to_id.csv", 0)
	negativeWords = ReadSourceReadonly("negative.txt", 0)
	positiveWords = ReadSourceReadonly("positive.txt", 0)

	classifier.Learn(positiveWords, Good)
	classifier.Learn(negativeWords, Bad)

	s = removeDuplicateStr(s)

	for i := range s {
		// stem
		s[i] = stemmer.Stem(s[i])

		// tokenize
		result := []string{}
		for _, word := range sastrawi.Tokenize(s[i]) {

			// check slang
			word = replaceSlang(word)

			if stopwords.Contains(word) || queryWords.Contains(word) {
				continue
			}
			result = append(result, word)
		}

		s[i] = strings.Join(result, " ")
	}

	mapPolarity := make(map[int]int)

	var positiveWords, negativeWords []string

	for i := range s {

		words := strings.Split(s[i], " ")
		for i := range words {
			words[i] = strings.TrimSpace(words[i])
		}

		scores, likely, _ := classifier.LogScores(words)
		fmt.Println(scores, " ", likely)

		if likely == 1 {
			positiveWords = append(positiveWords, words...)
		} else {
			negativeWords = append(negativeWords, words...)
		}

		mapPolarity[likely]++
	}

	fmt.Println("time passed : ", time.Since(t).Milliseconds(), " ms")

	for k, v := range mapPolarity {
		fmt.Println("key = ", k, " | ", v)
	}

	mapPositiveWords, mapNegativeWords := countPositiveAndNegativeWords(positiveWords, negativeWords)

	mostAppearancesPositive := getMostAppearancesWords(10, mapPositiveWords)
	mostAppearancesNegative := getMostAppearancesWords(10, mapNegativeWords)

	fmt.Println("mostAppearancesPositive: ", mostAppearancesPositive)
	fmt.Println("mostAppearancesNegative: ", mostAppearancesNegative)

	// 	// probs, likely, _ := classifier.ProbScores(
	// 	// 	strings.Split(s[i], " "),
	// 	// )

	// 	fmt.Println(scores, " ", likely)
	// }

	// scores, likely, _ := classifier.LogScores(
	// 	[]string{"terlalu murah",
	// 		"terlalu tinggi",
	// 		"terlambat"},
	// )
	// fmt.Println(scores, " ", likely)
}

func countPositiveAndNegativeWords(positiveWords, negativeWords []string) (map[string]int, map[string]int) {
	mapPositiveWords := make(map[string]int)
	mapNegativeWords := make(map[string]int)

	// breakdown positive and negative words
	for _, positiveWord := range positiveWords {
		if _, ok := mapPositiveWords[positiveWord]; ok {
			mapPositiveWords[positiveWord]++
			continue
		}
		mapPositiveWords[positiveWord]++
	}

	// breakdown negative and negative words
	for _, negativeWord := range negativeWords {
		if _, ok := mapPositiveWords[negativeWord]; ok {
			mapNegativeWords[negativeWord]++
			continue
		}
		mapNegativeWords[negativeWord]++
	}

	return mapPositiveWords, mapNegativeWords
}

func getMostAppearancesWords(rank int, mapWordFrequencies map[string]int) (mapMostAppearancesWord []string) {
	type WordAppearances struct {
		Key   string
		Value int
	}

	var wordFrequencies []WordAppearances
	for k, v := range mapWordFrequencies {
		wordFrequencies = append(wordFrequencies, WordAppearances{k, v})
	}

	sort.Slice(wordFrequencies, func(i, j int) bool {
		return wordFrequencies[i].Value > wordFrequencies[j].Value
	})

	for i := range wordFrequencies {
		if i == 10 {
			break
		}

		mapMostAppearancesWord = append(mapMostAppearancesWord, fmt.Sprintf("%s|%d", wordFrequencies[i].Key, wordFrequencies[i].Value))
	}

	return mapMostAppearancesWord
}
