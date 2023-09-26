package speedtest

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type ResultReader struct{}

const SOURCE_FILE = "result.csv"

func (r ResultReader) ReadAll() ([]DownloadTestResult, error) {
	file, err := os.Open(SOURCE_FILE)
	if err != nil {
		return []DownloadTestResult{}, fmt.Errorf("failed to open result.csv: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var results []DownloadTestResult

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error reading record %s", err)
		}

		if len(record) < 5 {
			return nil, fmt.Errorf("invalid record: %v", record)
		}

		downloadSpeed, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid download speed %s", record[3])
		}

		ping, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ping: %s", record[4])
		}

		result := DownloadTestResult{
			Timestamp:     record[0],
			HostIP:        record[1],
			NetworkName:   record[2],
			DownloadSpeed: downloadSpeed,
			Ping:          ping,
		}

		results = append(results, result)
	}
	return results, nil
}

func (r ResultReader) ReadLastN(n int) ([]DownloadTestResult, error) {
	if n <= 0 {
		return []DownloadTestResult{}, fmt.Errorf("n shoud be a postive integer: %d given", n)
	}

	allResults, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if n > len(allResults) {
		return allResults, nil
	}

	return allResults[len(allResults)-n:], nil
}
