package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	once      sync.Once
	jsonSlang string
	mSlang    = make(map[string]string)
)

func replaceSlang(s string) string {
	once.Do(initSlang)

	if slangResult, ok := mSlang[s]; ok {
		return slangResult
	}

	return s
}

func initSlang() {
	content, err := ioutil.ReadFile("dictionary/slang.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &mSlang)
	if err != nil {
		log.Fatal(err)
	}

}
