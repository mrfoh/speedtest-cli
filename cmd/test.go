package speedtest

import (
	"fmt"

	speedtest "github.com/mrfoh/speedtest/pkg"
	"github.com/spf13/cobra"
)

const SPEED_UNIT = "Mbps"
const PING_UNIT = "ms"

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test your internet speed",
	Run: func(cmd *cobra.Command, args []string) {
		resultMsg := "IP Address: %s\nISP Name: %s\nDownload speed: %.2f%s\nPing: %.2f%s"
		fmt.Println("Testing your internet speed...")
		downloader := speedtest.RealDownloader{}
		ipGetter := speedtest.RealHostIPGetter{}

		writer, err := speedtest.NewRealWriter("result.csv")
		if err != nil {
			fmt.Println(err)
			return
		}

		timeProvider := speedtest.RealTimeProvider{}

		test := speedtest.NewDownloadTest(&downloader, writer, &timeProvider, &ipGetter)
		result, err := test.Run()

		if err != nil {
			fmt.Println(err)
			return
		}

		resultString := fmt.Sprintf(resultMsg, result.HostIP, result.NetworkName, result.DownloadSpeed, SPEED_UNIT, result.Ping, PING_UNIT)

		fmt.Println(resultString)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
