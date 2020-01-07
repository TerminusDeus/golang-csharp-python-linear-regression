package main

import (
	"encoding/csv"
	"fmt"
	"github.com/xenolf/lego/log"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"os"
	"strconv"
)

// based on https://medium.com/devthoughts/linear-regression-with-go-ff1701455bcd

func main() {
	// we open the csv file from the disk:
	f, err := os.Open("kc_house_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// we create a new csv reader specifying the number of column it has:
	salesData := csv.NewReader(f)
	salesData.FieldsPerRecord = 21

	// we read all the records:
	records, err := salesData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// by slicing the records we skip the header:
	records = records[1:]

	// we iterate over all the records and keep track of all the gathered values
	// for each column:
	columnValues := map[int]plotter.Values{}
	for i, record := range records {
		// we want one histogram per column, so we will iterate over all the columns
		// we have and gather the data for each in separate value set in columnValues
		// we are skipping the ID column and the Date, so we start on index 2:
		for c := 2; c < salesData.FieldsPerRecord; c++ {
			if _, found := columnValues[c]; !found {
				columnValues[c] = make(plotter.Values, len(records))
			}
			// we parse each close value and add it to our set:
			floatVal, err := strconv.ParseFloat(record[c], 64)
			if err != nil {
				log.Fatal(err)
			}
			columnValues[c][i] = floatVal
		}
	}
	// once we have all the data, we draw each graph:
	for c, values := range columnValues {
		// create a new plot:
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("Histogram of %s", records[0][c])

		// create a new normalized histogram and add it to the plot:
		h, err := plotter.NewHist(values, 16)
		if err != nil {
			log.Fatal(err)
		}
		h.Normalize(1)
		p.Add(h)

		// save the plot to a PNG file:
		if err := p.Save(
			10*vg.Centimeter,
			10*vg.Centimeter,
			fmt.Sprintf("Histogram of %s", records[0][c]),
		); err != nil {
			log.Fatal(err)
		}
	}
}
