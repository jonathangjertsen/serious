package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available ports",
	Run: func(cmd *cobra.Command, args []string) {
		ports, err := serial.GetPortsList()

		if err != nil {
			log.Fatal(err)
		}

		for _, port := range ports {
			fmt.Printf("%v\n", port)
		}
	},
}
