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

type WordAppearances struct {
	Key   string
	Value int
}

type TwitterData struct {
	Tweet                  string
	Date                   string
	Sentiment              string
	TweetAfterPreprocessed string
}

var (
	stopwords           sastrawi.Dictionary
	dictionary          sastrawi.Dictionary
	queryWords          sastrawi.Dictionary
	additionalStopWords sastrawi.Dictionary
	stemmer             sastrawi.Stemmer

	NegativeWords []string
	PositiveWords []string

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
	queryWordsList = []string{"goto", "ipo", "gojek", "tokopedia", "tokped", "toped"}
)

func init() {
	additionalStopWordList := ReadSourceReadonly("dictionary/stopwords.txt", 0)

	stopwords = sastrawi.DefaultStopword()
	dictionary = sastrawi.DefaultDictionary()
	stemmer = sastrawi.NewStemmer(dictionary)
	queryWords = sastrawi.NewDictionary(queryWordsList...)
	additionalStopWords = sastrawi.NewDictionary(additionalStopWordList...)

	proceedDataTraining()

	// if 0 then Bad
	// if 1 then Good
	classifier = bayesian.NewClassifier(Bad, Good)

}

func main() {

	t := time.Now()

	// read tweet data
	tweetData := ReadTwitterData("tweets/test_data.csv")
	fmt.Println("Tweets: ", len(tweetData))

	tweetData = removeDuplicateTweet(tweetData)
	fmt.Println("Tweets removed duplicate: ", len(tweetData))

	NegativeWords = append(NegativeWords, ReadSourceReadonly("dictionary/negative.txt", 0)...)
	PositiveWords = append(PositiveWords, ReadSourceReadonly("dictionary/positive.txt", 0)...)

	// do sentiment analysis
	result, positiveWordsSentiment, negativeWordsSentiment := analyzeSentiment(tweetData, PositiveWords, NegativeWords)

	// write sentiment result
	writeSentimentResult("result/twit_result_sentiment.csv", result)

	fmt.Println("time passed : ", time.Since(t).Milliseconds(), " ms")

	// get mot appearance words
	mapPositiveWords, mapNegativeWords := countPositiveAndNegativeWords(positiveWordsSentiment, negativeWordsSentiment)
	writeMostAppearanceWords("result/most_appearances_positive.csv", getMostAppearancesWords(100, mapPositiveWords))
	writeMostAppearanceWords("result/most_appearances_negative.csv", getMostAppearancesWords(100, mapNegativeWords))

	// kfold cross validation analysis
	kFoldCrossValidation(10, result)

}

func proceedDataTraining() {
	trainingDataStr := ReadSourceReadonly("tweets/training_data.csv", 1)
	trainingDataSentiment := ReadSourceReadonly("tweets/training_data.csv", 0)

	for i := range trainingDataStr {
		trainingDataStr[i] = preProcessString(trainingDataStr[i])
	}

	for i := range trainingDataStr {
		if trainingDataSentiment[i] == "Positive" {
			words := strings.Split(trainingDataStr[i], " ")
			PositiveWords = append(PositiveWords, words...)
		} else if trainingDataSentiment[i] == "Negative" {
			words := strings.Split(trainingDataStr[i], " ")
			NegativeWords = append(NegativeWords, words...)
		}
	}
}

func removeDuplicateTweet(data []TwitterData) []TwitterData {
	mapStr := make(map[string]bool)
	listTweetData := []TwitterData{}

	for _, item := range data {
		if _, value := mapStr[item.Tweet]; !value {
			mapStr[item.Tweet] = true
			listTweetData = append(listTweetData, item)
		}
	}

	return listTweetData
}

func analyzeSentiment(tweetData []TwitterData, positiveWords, negativeWords []string) ([]TwitterData, []string, []string) {
	var (
		positiveWordsSentiment, negativeWordsSentiment []string
		positiveSentiment, negativeSentiment           int
	)

	classifier.Learn(positiveWords, Good)
	classifier.Learn(negativeWords, Bad)

	result := make([]TwitterData, len(tweetData))

	for i := range tweetData {
		result[i].Date = tweetData[i].Date
		result[i].Tweet = tweetData[i].Tweet
		result[i].TweetAfterPreprocessed = preProcessString(tweetData[i].Tweet)
	}

	for i := range tweetData {

		words := strings.Split(result[i].TweetAfterPreprocessed, " ")
		for i := range words {
			words[i] = strings.TrimSpace(words[i])
		}

		scores, likely, _ := classifier.LogScores(words)
		_ = scores

		if likely == 1 {
			positiveWordsSentiment = append(positiveWordsSentiment, words...)
			result[i].Sentiment = "Positive"
			positiveSentiment++
		} else {
			negativeWordsSentiment = append(negativeWordsSentiment, words...)
			result[i].Sentiment = "Negative"
			negativeSentiment++
		}
	}

	fmt.Println("Positive Sentiment: ", positiveSentiment)
	fmt.Println("Negative Sentiment: ", negativeSentiment)

	return result, positiveWordsSentiment, negativeWordsSentiment
}

func preProcessString(s string) string {
	result := []string{}

	// stem
	s = stemmer.Stem(s)

	// tokenize
	for _, word := range sastrawi.Tokenize(s) {

		// check slang
		word = replaceSlang(word)

		if stopwords.Contains(word) {
			continue
		}

		result = append(result, word)
	}

	s = strings.Join(result, " ")

	return s
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

func getMostAppearancesWords(rank int, mapWordFrequencies map[string]int) (mapMostAppearancesWord []WordAppearances) {

	var wordFrequencies []WordAppearances
	for k, v := range mapWordFrequencies {
		wordFrequencies = append(wordFrequencies, WordAppearances{k, v})
	}

	sort.Slice(wordFrequencies, func(i, j int) bool {
		return wordFrequencies[i].Value > wordFrequencies[j].Value
	})

	for i := range wordFrequencies {
		if i == rank {
			break
		}

		mapMostAppearancesWord = append(mapMostAppearancesWord, wordFrequencies[i])
	}

	return mapMostAppearancesWord
}
