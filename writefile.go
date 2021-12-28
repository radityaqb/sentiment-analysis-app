package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func writeSentimentResult(filename string, tweetData []TwitterData) {
	csvResult := [][]string{}

	for i := range tweetData {
		r := []string{tweetData[i].Date, tweetData[i].Sentiment, tweetData[i].Tweet, tweetData[i].TweetAfterPreprocessed}
		csvResult = append(csvResult, r)
	}

	csvFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	for _, row := range csvResult {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
}

func writeMostAppearanceWords(filename string, mostAppearancesWord []WordAppearances) {
	csvResult := [][]string{}

	csvFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	for i := range mostAppearancesWord {
		r := []string{mostAppearancesWord[i].Key, fmt.Sprint(mostAppearancesWord[i].Value)}
		csvResult = append(csvResult, r)
	}

	for _, row := range csvResult {
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
}
