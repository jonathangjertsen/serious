package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	hw "github.com/jonathangjertsen/serious/hw"
	messages "github.com/jonathangjertsen/serious/messages"
	ui "github.com/jonathangjertsen/serious/ui"
	"github.com/spf13/cobra"
)

const version = "v0.0.0"

type UsageError struct {
	Opt     string
	Value   string
	Allowed []string
}

func (e *UsageError) Error() string {
	allowed := ""
	if len(e.Allowed) > 0 {
		allowed = fmt.Sprintf(" (allowed: %s)", strings.Join(e.Allowed, ", "))
	}
	return fmt.Sprintf("%s can not be '%s'%s", e.Opt, e.Value, allowed)
}

var RootCmd = &cobra.Command{
	Use:   "serious",
	Short: fmt.Sprintf("serious serial CLI %s", version),
	Long: fmt.Sprintf(`serious serial CLI %s

When run with no parameters, serious runs in interactive mode.
`, version),
	Run: func(cmd *cobra.Command, args []string) {
		logfile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logfile)

		uiStr, err := cmd.Flags().GetString("ui")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		uiImpl, err := GetUi(uiStr)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		channel := make(chan messages.Message)

		var wg sync.WaitGroup
		wg.Add(1)
		go uiWorker(uiImpl, &channel, &wg)
		wg.Add(1)
		go hwWorker(&channel, &wg)

		wg.Wait()
	},
	Version:           version,
	DisableAutoGenTag: true,
}

func init() {
	RootCmd.PersistentFlags().String("ui", "terminal", "UI types")
}

func Execute() {
	cobra.MousetrapHelpText = ""
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func GetUi(uiStr string) (ui.Ui, error) {
	if uiStr == "terminal" {
		return ui.NewTerminal(), nil
	}
	return nil, fmt.Errorf("--ui can not be '%s' (allowed values: 'terminal' [default])", uiStr)
}

func uiWorker(uiImpl ui.Ui, channel *chan messages.Message, wg *sync.WaitGroup) {
	defer logPanic()
	defer wg.Done()
	log.Printf("UI worker started")
	uiImpl.StartReadTask(time.Millisecond * 250)
	uiImpl.Run(channel)
	log.Printf("UI worker exited")
}

func hwWorker(channel *chan messages.Message, wg *sync.WaitGroup) {
	defer logPanic()
	defer wg.Done()
	hwImpl, err := hw.NewSerial(channel)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	hwImpl.Run()
}

func logPanic() {
	if p := recover(); p != nil {
		log.Printf("Panic: %+v", string(debug.Stack()))
		panic(p)
	}
}
