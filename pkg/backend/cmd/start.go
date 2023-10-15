package cmd

import (
	"github.com/doornoc/dsbd-ssh/pkg/api/core/tool"
	"github.com/spf13/cobra"
	"log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start controller",
	Long:  ``,
}

var startOnceCmd = &cobra.Command{
	Use:   "once",
	Short: "start for once",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//err = remote.GettingDeviceConfig(true)
		//if err != nil {
		//	notify.NotifyErrorToSlack(err)
		//}

		log.Println("end")
	},
}

var startManualCmd = &cobra.Command{
	Use:   "manual",
	Short: "start for cron",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error = nil
		tool.Debug, err = cmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		log.Println("end")

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(startOnceCmd)
	startCmd.AddCommand(startManualCmd)
	startCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode")
}
