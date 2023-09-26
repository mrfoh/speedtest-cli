package speedtest

import (
	"fmt"

	speedtest "github.com/mrfoh/speedtest/pkg"
	"github.com/spf13/cobra"
)

var averageCmd = &cobra.Command{
	Use:   "average",
	Short: "Show average of previous test results",
	Run: func(cmd *cobra.Command, args []string) {
		getAllResults, _ := cmd.Flags().GetBool("all")
		getLastNResults, _ := cmd.Flags().GetInt("last")
		var results []speedtest.DownloadTestResult = []speedtest.DownloadTestResult{}

		reader := speedtest.ResultReader{}

		if getAllResults {
			results, _ = reader.ReadAll()
		}

		if getLastNResults > 0 {
			results, _ = reader.ReadLastN(getLastNResults)
		}

		resultCount := len(results)
		var totalDownloadSpeed float64 = 0
		var totalPing float64 = 0

		for _, result := range results {
			totalDownloadSpeed += result.DownloadSpeed
			totalPing += result.Ping
		}

		averageDownloadSpeed := totalDownloadSpeed / float64(resultCount)
		averagePing := totalPing / float64(resultCount)

		averageLn := fmt.Sprintf("Average download speed: %.2fMbps\nAverage Ping: %.2fms\n", averageDownloadSpeed, averagePing)

		fmt.Print(averageLn)
	},
}

func init() {
	rootCmd.AddCommand(averageCmd)
	averageCmd.Flags().BoolP("all", "a", true, "Calculate average of all results")
	averageCmd.Flags().IntP("last", "l", 0, "Calculate average of last n results")
}
