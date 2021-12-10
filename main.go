package main

import (
	"fmt"
	"strings"
	"time"

	sastrawi "github.com/RadhiFadlillah/go-sastrawi"
	"github.com/jbrukh/bayesian"
)

var (
	stopwords  sastrawi.Dictionary
	dictionary sastrawi.Dictionary
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

func init() {

	stopwords = sastrawi.DefaultStopword()
	dictionary = sastrawi.DefaultDictionary()
	stemmer = sastrawi.NewStemmer(dictionary)

	// if 0 then Bad
	// if 1 then Good
	classifier = bayesian.NewClassifier(Bad, Good)

}

func main() {
	t := time.Now()

	fmt.Println("Sentiment Analysis App")

	// 1. Read
	s = ReadSourceReadonly("twit_translated_to_id.csv", 1)
	negativeWords = ReadSourceReadonly("negative.txt", 0)
	positiveWords = ReadSourceReadonly("positive.txt", 0)

	// // remove newline
	// re := regexp.MustCompile(`\r?\n`)
	// for i := range s {
	// 	s[i] = re.ReplaceAllString(s[i], " ")
	// 	fmt.Println(s[i])
	// }
	// log.Fatal("done")

	classifier.Learn(positiveWords, Good)
	classifier.Learn(negativeWords, Bad)

	for i := range s {

		// tokenize
		result := []string{}
		for _, word := range sastrawi.Tokenize(s[i]) {
			if stopwords.Contains(word) {
				continue
			}
			result = append(result, stemmer.Stem(word))
		}

		s[i] = strings.Join(result, " ")
	}

	s = removeDuplicateStr(s)
	// fmt.Println(len(s))

	// fmt.Println(s[1])
	mapPolarity := make(map[int]int)
	for i := range s {
		fmt.Println(s[i])

		words := strings.Split(s[i], " ")
		for i := range words {
			words[i] = strings.TrimSpace(words[i])
		}

		scores, likely, _ := classifier.LogScores(words)
		fmt.Println(scores, " ", likely)
		mapPolarity[likely]++
	}

	fmt.Println("time passed : ", time.Since(t).Milliseconds(), " ms")

	for k, v := range mapPolarity {
		fmt.Println("key = ", k, " | ", v)
	}

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
