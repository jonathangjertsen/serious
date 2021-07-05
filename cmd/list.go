package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	hw "github.com/jonathangjertsen/serious/hw"
	messages "github.com/jonathangjertsen/serious/messages"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available ports",
	Run: func(cmd *cobra.Command, args []string) {
		channel := make(chan messages.Message)

		var wg sync.WaitGroup
		wg.Add(1)
		go listHwWorker(&channel, &wg)

		wg.Add(1)
		go listDisplayWorker(&channel, &wg)

		wg.Wait()
	},
}

func listDisplayWorker(channel *chan messages.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	response := messages.SyncGetPorts(channel)
	for _, port := range response.Ports {
		selectionIndicator := ""
		if port == *response.OpenName {
			selectionIndicator = " [auto-selected]"
		}
		fmt.Printf("%v%s\n", port, selectionIndicator)
	}
	messages.SyncExit(channel)
}

func listHwWorker(channel *chan messages.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	hwImpl, err := hw.NewSerial(channel)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	hwImpl.Run()
}
