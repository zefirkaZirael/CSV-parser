package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var csvparser CSVParser = &CSVParserImpl{}

	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}

		fmt.Println("line:", line, "[", csvparser.GetNumberOfFields(), "]")
		field, err := csvparser.GetField(0)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("field:", field)
	}
}
