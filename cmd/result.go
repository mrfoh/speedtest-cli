package speedtest

import (
	"fmt"

	speedtest "github.com/mrfoh/speedtest/pkg"
	"github.com/spf13/cobra"
)

var resultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Show previous test results",
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

		for _, result := range results {
			resultLn := fmt.Sprintf("%s - %s - %s; %.2fMbps - %.2fms\n", result.Timestamp, result.HostIP, result.NetworkName, result.DownloadSpeed, result.Ping)
			fmt.Print(resultLn)
		}
	},
}

func init() {
	rootCmd.AddCommand(resultsCmd)
	resultsCmd.Flags().BoolP("all", "a", true, "Show all results")
	resultsCmd.Flags().IntP("last", "l", 0, "Show last n results")
}
