package cmd

import (
	"log"

	"github.com/mfdeux/snoo-sockets/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("host", "h", "127.0.0.1", "host to use")
	serveCmd.Flags().IntP("port", "p", 8900, "port to use")
}

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		err := server.NewWebsocketServer(host, port)
		if err != nil {
			log.Fatal(err)
		}
	},
}
