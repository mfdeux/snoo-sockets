package cmd

import (
	"fmt"
	"log"

	"github.com/mfdeux/snoo-sockets/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(consumeCmd)
	consumeCmd.Flags().StringP("endpoint", "e", "ws://127.0.0.1:8200", "endpoint of websocket server")
	consumeCmd.Flags().StringP("output", "o", "stdout", "output of stream")
}

var consumeCmd = &cobra.Command{
	Use: "consume",
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, _ := cmd.Flags().GetString("endpoint")
		snoo, err := client.NewClient(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		messages := snoo.Consume()
		for message := range messages {
			fmt.Println(string(message.([]byte)))
		}
	},
}
