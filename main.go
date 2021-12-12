package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
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

	// data test, scrapped from twitter
	s      []string
	sRaw   []string
	date   []string
	result []int

	classifier *bayesian.Classifier
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
	s = ReadSourceReadonly("twit_translated_to_id_date.csv", 1)
	date = ReadSourceReadonly("twit_translated_to_id_date.csv", 0)

	// clone sRaw from s
	for i := range s {
		sRaw = append(sRaw, s[i])
	}

	result = make([]int, len(s))

	negativeWords = ReadSourceReadonly("negative_mk2.txt", 0)
	positiveWords = ReadSourceReadonly("positive_mk2.txt", 0)

	classifier.Learn(positiveWords, Good)
	classifier.Learn(negativeWords, Bad)

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
		_ = scores
		// fmt.Println(scores, " ", likely)

		if likely == 1 {
			positiveWords = append(positiveWords, words...)
		} else {
			negativeWords = append(negativeWords, words...)
		}

		mapPolarity[likely]++

		result[i] = likely
	}

	fmt.Println("time passed : ", time.Since(t).Milliseconds(), " ms")

	csvResult := [][]string{}

	for i := range s {
		resultStr := "Negative"
		if result[i] > 0 {
			resultStr = "Positive"
		}

		// fmt.Printf("%s | %s | %s\n", date[i], resultStr, sRaw[i])
		r := []string{date[i], resultStr, sRaw[i]}
		csvResult = append(csvResult, r)
	}

	//// Write to CSV
	csvFile, err := os.Create("twit_result_sentiment.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	for _, row := range csvResult {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()

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
