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

type Downloader interface {
	Download(url string) (response *http.Response, err error)
}

type HostIPGetter interface {
	GetHostIP() (string, error)
}

type Writer interface {
	Write(record []string) error
}

type TimeProvider interface {
	Now() time.Time
}

type RealDownloader struct{}

type RealHostIPGetter struct{}

func (rd *RealDownloader) Download(url string) (*http.Response, error) {
	return http.Get(url)
}

type RealWriter struct {
	file *os.File
}

func NewRealWriter(filename string) (*RealWriter, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &RealWriter{file: file}, nil
}

func (rw *RealWriter) Write(record []string) error {
	writer := csv.NewWriter(rw.file)
	defer writer.Flush()
	return writer.Write(record)
}

type RealTimeProvider struct{}

func (rtp *RealTimeProvider) Now() time.Time {
	return time.Now()
}

type DownloadTest struct {
	downloader   Downloader
	writer       Writer
	timeProvider TimeProvider
	hostIPGetter HostIPGetter
}

type DownloadTestResult struct {
	Timestamp     string
	HostIP        string
	NetworkName   string
	DownloadSpeed float64
	Ping          float64
}

func NewDownloader() *RealDownloader {
	return &RealDownloader{}
}

func NewDownloadTest(downloader Downloader, writer Writer, timeProvider TimeProvider, ipGetter HostIPGetter) *DownloadTest {
	return &DownloadTest{
		downloader:   downloader,
		writer:       writer,
		timeProvider: timeProvider,
		hostIPGetter: ipGetter,
	}
}

func (d *RealHostIPGetter) GetHostIP() (string, error) {
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

func (d *DownloadTest) WriteResultToFile(result DownloadTestResult) error {
	record := []string{
		result.Timestamp,
		result.HostIP,
		result.NetworkName,
		fmt.Sprintf("%.2f", result.DownloadSpeed),
		fmt.Sprintf("%.2f", result.Ping),
	}
	return d.writer.Write(record)
}

func (d *DownloadTest) Run() (DownloadTestResult, error) {
	TEST_FILE_URL := "https://ispindex.s3.eu-west-2.amazonaws.com/testsfiles/10MB.zip"
	timestamp := d.timeProvider.Now().Format(time.RFC3339)
	hostIP, _ := d.hostIPGetter.GetHostIP()
	result := DownloadTestResult{
		Timestamp:   timestamp,
		HostIP:      hostIP,
		NetworkName: "Unknown Network",
	}
	startTime := d.timeProvider.Now()

	response, err := d.downloader.Download(TEST_FILE_URL)
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

	err = d.WriteResultToFile(result)
	if err != nil {
		return result, err
	}

	return result, nil
}
