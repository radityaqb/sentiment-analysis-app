package main

import (
	"fmt"
	"strings"
)

func kFoldCrossValidation(fold int, tweetData []TwitterData) {

	testDataCount := len(tweetData) / fold

	var lastIndexSaved int

	for i := 0; i < fold; i++ {
		fmt.Println("================================")
		fmt.Println("Fold ", i)

		var (
			firstIndex, lastIndex        int
			testData, trainingData       []TwitterData
			positiveWords, negativeWords []string
		)

		if i == 0 {
			firstIndex = 0
		} else {
			firstIndex = lastIndexSaved + 1
		}

		lastIndex = firstIndex + (testDataCount - 1)
		lastIndexSaved = lastIndex

		if i == fold {
			remainder := len(tweetData) % fold
			lastIndex += remainder
		}

		// fmt.Println("First Index: ", firstIndex)
		// fmt.Println("Last Index: ", lastIndex)

		// breakdown test data and training data
		for j := range tweetData {
			if j >= firstIndex && j <= lastIndex {
				testData = append(testData, tweetData[j])
			} else {
				trainingData = append(trainingData, tweetData[j])
			}
		}

		// get positive and negative words from test data
		for i := range testData {
			if testData[i].Sentiment == "Positive" {
				words := strings.Split(testData[i].TweetAfterPreprocessed, " ")
				positiveWords = append(PositiveWords, words...)
			} else if testData[i].Sentiment == "Negative" {
				words := strings.Split(testData[i].TweetAfterPreprocessed, " ")
				negativeWords = append(NegativeWords, words...)
			}
		}

		// fmt.Println("Test Data: ", len(testData))
		// fmt.Println("Training Data: ", len(trainingData))

		// do sentiment analysis
		result, _, _ := analyzeSentiment(trainingData, positiveWords, negativeWords)

		var (
			accuracy, precision, recall                              float64
			truePositive, trueNegative, falsePositive, falseNegative float64
		)

		for i := range result {
			// if i == 10 {
			// 	fmt.Println("Training Data ", trainingData[i].TweetAfterPreprocessed, "Sentiment: ", trainingData[i].Sentiment)
			// 	fmt.Println("Result Data ", result[i].TweetAfterPreprocessed, "Sentiment: ", result[i].Sentiment)
			// }

			if result[i].Sentiment == "Positive" {
				if result[i].Sentiment == trainingData[i].Sentiment {
					truePositive++
				} else if result[i].Sentiment != trainingData[i].Sentiment {
					falsePositive++
				}
			}

			if result[i].Sentiment == "Negative" {
				if result[i].Sentiment == trainingData[i].Sentiment {
					trueNegative++
				} else if result[i].Sentiment != trainingData[i].Sentiment {
					falseNegative++
				}
			}
		}

		fmt.Println("True Positive: ", truePositive)
		fmt.Println("False Positive: ", falsePositive)
		fmt.Println("True Negative: ", trueNegative)
		fmt.Println("False Negative: ", falseNegative)

		accuracy = (truePositive + trueNegative) / float64(len(result)) * 100
		precision = truePositive / (truePositive + falsePositive) * 100
		recall = truePositive / (truePositive + falseNegative) * 100

		fmt.Println("Accuracy: ", accuracy, " ")
		fmt.Println("Precission: ", precision)
		fmt.Println("Recall: ", recall)
	}
}
