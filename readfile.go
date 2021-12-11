package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ReadSourceReadonly(filename string, idx int) []string {
	f, err := os.Open(fmt.Sprint(filename))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		// log.Fatal("Unable to parse file as CSV for "+filename, err)
		// id,permalink,username,text,date,polarity,pos_w,neu_w,neg_w
		fmt.Println("Unable to parse file as CSV for "+filename, err)
	}

	result := []string{}
	for i := range records {

		// skip header
		if i == 0 {
			continue
		}
		result = append(result, records[i][idx])
	}

	return result
}

func ReadSourceJSON(filename string) map[string]string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	mapWords := make(map[string]string)

	err = json.Unmarshal(byteValue, &mapWords)
	if err != nil {
		fmt.Println(err)
	}

	return mapWords
}
