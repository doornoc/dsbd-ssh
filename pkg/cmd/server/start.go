package server

import (
	"github.com/doornoc/dsbd-ssh/pkg/api/core"
	"github.com/doornoc/dsbd-ssh/pkg/api/core/tool"
	"github.com/spf13/cobra"
	"log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start grpc server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		tool.Port = uint(port)

		core.Server()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().IntP("port", "p", 50051, "grpc port")
}
