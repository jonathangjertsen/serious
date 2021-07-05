package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	fmt.Print(portListString(&response))
	messages.SyncExit(channel)
}

func portListString(response *messages.PortsResponse) string {
	var sb strings.Builder
	for _, port := range response.Ports {
		selectionIndicator := ""
		if port == response.OpenName {
			selectionIndicator = " [will be auto-selected]"
		}
		sb.WriteString(fmt.Sprintf("%v%s\n", port, selectionIndicator))
	}
	return sb.String()
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
