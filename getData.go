package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/csv"
	"encoding/json"
)

type Inner struct {
	Name			string			`json:"SC_FULLNM"`
	SCID			string			`json:"DISPID"`
	Price			FlexInt			`json:"pricecurrent"`
	PriceChange		FlexInt			`json:"pricechange"`
	Percentage		FlexInt			`json:"pricepercentchange"`
	High			FlexInt			`json:"HP,omitempty"`
	Low				FlexInt			`json:"LP,omitempty"`
	Volume			int				`json:"VOL"`
	LastUpdated		time.Time		`json:"lastupd"`
	LCL				FlexInt			`json:"lower_circuit_limit,omitempty"`
	UCL				FlexInt			`json:"upper_circuit_limit,omitempty"`
}

type Collection struct {
	Status			int				`json:"code"`
	Data			Inner			`json:"data"`
	URL				string			`json:"url"`
	LastScraped		time.Time		`json:"lastscraped"`
}

type FlexInt string

func (fi *FlexInt) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
			b = b[1 : len(b)-1]
	}
	i, err := strconv.ParseFloat(string(b), 32)
	*fi = FlexInt(fmt.Sprintf("%.2f", i))
	return
}

func createCSV(file *File) {
	var count int = 0
	var record Collection
	var dtlayout string = "02-01-2006 3:04:05 PM"

	csvFile, err := os.Create(file.directory + file.name)
	if err != nil {
		fmt.Printf("OS Error: File Creation Failed, %s\n", file.directory + file.name)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	writer.Write([]string{
		"Sno.",
		"Name",
		"SCID",
		"Current Price",
		"Change in Price",
		"Change in Price %",
		"High",
		"Low",
		"Volume",
		"Lower Control Limit",
		"Upper Control Limit",
		"Last Updated (in MoneyControl)",
		"Last Scraper",
	})

	files, err := ioutil.ReadDir(file.directory)
	if err != nil {
		fmt.Printf("IO Error: Error at Reading the Directory. %v\n", err)
	}

	for _, temp_file := range files {
		if strings.HasSuffix(temp_file.Name(), ".json") {
			var row []string
			count++

			jsonFile, err := os.Open(file.directory + temp_file.Name())
			if err != nil {
				fmt.Printf("IO Error: Error at Accessing the JSON file %s%s. %v\n", file.directory, temp_file.Name(), err)
			}
			defer jsonFile.Close()
	
			byteData, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				fmt.Printf("IO Error: Error at Reading the JSON file %s%s. %v\n", file.directory, temp_file.Name(), err)
			}
	
			err = json.Unmarshal(byteData, &record)
			if err != nil {
				fmt.Printf("JSON Error: Error at Unmarshal of the JSON file %s%s. %v\n", file.directory, temp_file.Name(), err)
			}
			row = append(row, strconv.Itoa(count))
			row = append(row, record.Data.Name)
			row = append(row, record.Data.SCID)
			row = append(row, fmt.Sprintf("%v", record.Data.Price))
			row = append(row, fmt.Sprintf("%v", record.Data.PriceChange))
			row = append(row, fmt.Sprintf("%v", record.Data.Percentage))
			row = append(row, fmt.Sprintf("%v", record.Data.High))
			row = append(row, fmt.Sprintf("%v", record.Data.Low))
			row = append(row, strconv.Itoa(record.Data.Volume))
			row = append(row, fmt.Sprintf("%v", record.Data.LCL))
			row = append(row, fmt.Sprintf("%v", record.Data.UCL))
			row = append(row, record.Data.LastUpdated.Format(dtlayout))
			row = append(row, record.LastScraped.Format(dtlayout))

			writer.Write(row)
		}
	}
	writer.Flush()
	fmt.Printf("\r%s\t\tCreated\t\t\t%v\n", file.directory + file.name, time.Now())
}

func deleteCSV(file *File) {
	err := os.Remove(file.directory + file.name)
    if err != nil {
        fmt.Printf("OS Error: File %s Removal failed with error: %v", file.name, err)
    }
    fmt.Printf("\r%s\t\tDeleted\t\t\t%v\n", file.directory + file.name, time.Now())
}
