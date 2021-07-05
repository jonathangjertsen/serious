package cmd

import (
	"fmt"
	"log"

	hw "github.com/jonathangjertsen/serious/hw"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available ports",
	Run: func(cmd *cobra.Command, args []string) {
		ser, err := hw.NewSerial()

		if err != nil {
			log.Fatal(err)
		}

		selectedIndex, _ := ser.Selected()
		for index, port := range ser.GetPorts() {
			selectionIndicator := ""
			if index == selectedIndex {
				selectionIndicator = " [auto-selected]"
			}
			fmt.Printf("%d %v%s\n", index, port, selectionIndicator)
		}
	},
}
