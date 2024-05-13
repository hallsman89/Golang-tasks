package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Data struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"location"`
}

func parseCsvFile() ([]Data, error) {
	var res []Data
	file, err := os.Open(csvFileName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	read := csv.NewReader(file)
	read.Comma = '\t'
	_, _ = read.Read()
	for {
		record, err := read.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data, err := makeData(record)
		if err != nil {
			return nil, err
		}
		res = append(res, data)
	}
	return res, nil
}

func makeData(record []string) (Data, error) {
	if len(record) != 6 {
		return Data{}, fmt.Errorf("Invalid person slice: %v", record)
	}
	id := record[0]
	name := record[1]
	address := record[2]
	phone := record[3]
	lon, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return Data{}, fmt.Errorf("Invalid Longitude: %s", record[4])
	}
	lat, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return Data{}, fmt.Errorf("Invalid Latitude: %s", record[5])
	}
	return Data{
		ID:      id,
		Name:    name,
		Address: address,
		Phone:   phone,
		Location: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{Lon: lon, Lat: lat},
	}, nil
}
