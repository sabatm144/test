package controller

import (
	"encoding/csv"
	"fmt"
	"os"
)

func LoadCsv() {

	// Open CSV file
	f, err := os.Open("controllers/load.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	// Loop through lines & turn into object
	for _, line := range lines {
		fmt.Println(line[0] + " " + line[1])
	}
}
