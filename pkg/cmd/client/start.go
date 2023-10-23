package client

import (
	"github.com/doornoc/dsbd-ssh/pkg/api/core"
	"github.com/spf13/cobra"
	"log"
)

// clientCmd represents the start command
var clientCmd = &cobra.Command{
	Use:   "start",
	Short: "start grpc client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		hostname, err := cmd.Flags().GetString("hostname")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		err = core.Client(hostname, port, username)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringP("hostname", "o", "localhost", "hostname")
	clientCmd.PersistentFlags().IntP("port", "p", 22, "port")
	clientCmd.PersistentFlags().StringP("username", "u", "", "username")
}
