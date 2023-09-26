package speedtest

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const TEST_FILE_URL = "https://ispindex.s3.eu-west-2.amazonaws.com/testsfiles/10MB.zip"

type DownloadTest struct{}

type DownloadTestResult struct {
	Timestamp     string
	HostIP        string
	NetworkName   string
	DownloadSpeed float64
	Ping          float64
}

func GetHostIP() (string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}

func WriteResultToFile(result DownloadTestResult) error {
	file, err := os.OpenFile("result.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open result.csv: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		result.Timestamp,
		result.HostIP,
		result.NetworkName,
		fmt.Sprintf("%.2f", result.DownloadSpeed),
		fmt.Sprintf("%.2f", result.Ping),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record to result.csv: %s", err)
	}

	return nil
}

func (d *DownloadTest) Run() (DownloadTestResult, error) {
	timestamp := time.Now().Format(time.RFC3339)
	hostIP, _ := GetHostIP()
	result := DownloadTestResult{
		Timestamp:   timestamp,
		HostIP:      hostIP,
		NetworkName: "Unknown Network",
	}
	startTime := time.Now()

	response, err := http.Get(TEST_FILE_URL)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		return result, err
	}

	elaspedTime := time.Since(startTime).Seconds()
	ping := time.Since(startTime).Milliseconds()
	fileSizeInBytes := float64(response.ContentLength)

	speedInBytesPerSecond := fileSizeInBytes / elaspedTime

	speedMbps := (speedInBytesPerSecond * 8) / 1_000_000

	result.DownloadSpeed = speedMbps
	result.Ping = float64(ping)

	err = WriteResultToFile(result)
	if err != nil {
		return result, err
	}

	return result, nil
}
