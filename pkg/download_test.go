package speedtest

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDownloader struct {
	mock.Mock
}

type MockHostIPGetter struct {
	mock.Mock
}

func (m *MockDownloader) Download(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(record []string) error {
	args := m.Called(record)
	return args.Error(0)
}

type MockTimeProvider struct {
	mock.Mock
}

func (m *MockTimeProvider) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockHostIPGetter) GetHostIP() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestRun(t *testing.T) {
	mockDownloader := new(MockDownloader)
	mockWriter := new(MockWriter)
	mockTimeProvider := new(MockTimeProvider)
	mockHostIpGetter := new(MockHostIPGetter)

	downloadTest := NewDownloadTest(mockDownloader, mockWriter, mockTimeProvider, mockHostIpGetter)

	// Mocking the time
	now := time.Now()
	mockTimeProvider.On("Now").Return(now)

	mockHostIpGetter.On("GetHostIP").Return("127.0.0.1", nil)

	// Mocking the downloader
	body := bytes.NewBufferString("fake http response body")
	mockDownloader.On("Download", mock.Anything).Return(&http.Response{
		Body:          io.NopCloser(body),
		ContentLength: int64(len(body.Bytes())),
	}, nil)

	// Mocking the writer
	mockWriter.On("Write", mock.Anything).Return(nil)

	result, err := downloadTest.Run()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDownloader.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
	mockTimeProvider.AssertExpectations(t)

	// You can also add more assertions to check the values in the result
	assert.Equal(t, now.Format(time.RFC3339), result.Timestamp)
	assert.Equal(t, "Unknown Network", result.NetworkName)
	assert.Equal(t, "127.0.0.1", result.HostIP)
}

func TestRun_Error(t *testing.T) {
	mockDownloader := new(MockDownloader)
	mockWriter := new(MockWriter)
	mockTimeProvider := new(MockTimeProvider)
	mockHostIpGetter := new(MockHostIPGetter)

	downloadTest := NewDownloadTest(mockDownloader, mockWriter, mockTimeProvider, mockHostIpGetter)

	// Mocking the time
	now := time.Now()
	mockTimeProvider.On("Now").Return(now)

	mockHostIpGetter.On("GetHostIP").Return("", errors.New("host IP error"))

	// Mocking the downloader to return an error
	mockDownloader.On("Download", mock.Anything).Return((*http.Response)(nil), errors.New("download error"))

	result, err := downloadTest.Run()

	assert.Error(t, err)
	assert.Equal(t, "download error", err.Error())
	assert.NotNil(t, result)
	mockDownloader.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
	mockTimeProvider.AssertExpectations(t)
}
