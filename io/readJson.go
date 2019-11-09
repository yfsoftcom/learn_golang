package main

import (
	"encoding/json"
	"os"
)

// define Type
type Feed struct {
	Name string `json:"site"`
	Link string `json:"link"`
	Type string `json:"type"`
}

func read(filename string) ([]*Feed, error) {
	// read file
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// close file when function close
	defer file.Close()

	// define a slice with the Feed's pointer
	var feeds []*Feed

	// read the filestream and convert to json
	err = json.NewDecoder(file).Decode(&feeds)

	// err 由
	return feeds, err
}

func main() {
	feeds, err := read("./demo.json")
	if err != nil {
		println("err:", err)
	}
	for _, feed := range feeds {
		println("Name:", feed.Name)
	}
}
