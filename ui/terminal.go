package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	messages "github.com/jonathangjertsen/serious/messages"
	"github.com/rivo/tview"
)

const (
	colorNone  = "[-]"
	colorDebug = "[grey]"
	colorError = "[red]"
	colorTx    = "[yellow:]"
	colorRx    = "[blue:]"
)

type Terminal struct {
	connected      bool
	channel        *chan messages.Message
	app            *tview.Application
	widgets        []tview.Primitive
	output         *tview.TextView
	baudSelect     *tview.InputField
	deviceSelect   *tview.DropDown
	dataBitsSelect *tview.DropDown
	stopBitsSelect *tview.DropDown
	paritySelect   *tview.DropDown
}

func NewTerminal() *Terminal {
	// Get default port config
	cfg := messages.DefaultPortConfig()

	// Init terminal
	term := &Terminal{}

	// Init tview app
	term.app = tview.NewApplication()

	// Add main layout object
	// Each item that gets added gets a new row
	main := tview.NewFlex().SetDirection(tview.FlexRow)
	term.app.SetRoot(main, true)

	// Add header box at the start with 2 columns
	headerBox := tview.NewFlex()
	main.AddItem(headerBox, 6, 1, false)

	// Row to the right
	headerConfigItemRows := tview.NewFlex().SetDirection(tview.FlexRow)
	headerBox.AddItem(headerConfigItemRows, 0, 3, false)

	// First row
	headerConfigItemRow0 := tview.NewFlex()
	headerConfigItemRows.AddItem(headerConfigItemRow0, 0, 1, false)

	headerConfigItemRow1 := tview.NewFlex()
	headerConfigItemRows.AddItem(headerConfigItemRow1, 0, 1, false)

	// Add device config
	device := tview.NewDropDown()
	device.SetLabel("device ")
	device.SetFieldWidth(20)
	device.SetBorder(true)
	device.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	device.AddOption("None", nil)
	device.SetCurrentOption(0)
	term.widgets = append(term.widgets, device)
	term.deviceSelect = device
	headerConfigItemRow0.AddItem(device, 0, 2, false)

	// Add baud rate config
	baud := tview.NewInputField()
	baud.SetLabel("baud ")
	baud.SetFieldWidth(8)
	baud.SetAcceptanceFunc(tview.InputFieldInteger)
	baud.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	baud.SetText(strconv.Itoa(cfg.BaudRate))
	baud.SetBorder(true)
	term.widgets = append(term.widgets, baud)
	term.baudSelect = baud
	headerConfigItemRow0.AddItem(baud, 0, 2, false)

	// Add data bits config
	dataBits := tview.NewDropDown()
	dataBits.SetLabel("bits ")
	dataBits.SetFieldWidth(6)
	dataBits.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	dataBits.SetBorder(true)
	dataBits.AddOption("5", nil)
	dataBits.AddOption("6", nil)
	dataBits.AddOption("7", nil)
	dataBits.AddOption("8", nil)
	dataBits.SetCurrentOption(3)
	term.widgets = append(term.widgets, dataBits)
	term.dataBitsSelect = dataBits
	headerConfigItemRow1.AddItem(dataBits, 0, 1, false)

	// Add parity config
	parity := tview.NewDropDown()
	parity.SetLabel("parity ")
	parity.SetFieldWidth(6)
	parity.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	parity.SetBorder(true)
	parity.AddOption(messages.ParityNone, nil)
	parity.AddOption(messages.ParityOdd, nil)
	parity.AddOption(messages.ParityEven, nil)
	parity.AddOption(messages.ParityAlways1, nil)
	parity.AddOption(messages.ParityAlways0, nil)
	parity.SetCurrentOption(0)
	term.widgets = append(term.widgets, parity)
	term.paritySelect = parity
	headerConfigItemRow1.AddItem(parity, 0, 1, false)

	// Add stop bits config
	stopBits := tview.NewDropDown()
	stopBits.SetLabel("stop ")
	stopBits.SetFieldWidth(6)
	stopBits.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	stopBits.SetBorder(true)
	stopBits.AddOption(messages.StopBitsOne, nil)
	stopBits.AddOption(messages.StopBitsOnePointFive, nil)
	stopBits.AddOption(messages.StopBitsTwo, nil)
	stopBits.SetCurrentOption(0)
	term.widgets = append(term.widgets, stopBits)
	term.stopBitsSelect = stopBits
	headerConfigItemRow1.AddItem(stopBits, 0, 1, false)

	// Add terminator select button
	// Need to forward declare input field since this dropdown affects placeholder value
	var input *tview.InputField
	entry := tview.NewDropDown()
	entry.SetLabel("entry ")
	entry.SetFieldWidth(14)
	entry.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
		_, entryStr := entry.GetCurrentOption()
		if entryStr == "Immediate" {
			input.SetPlaceholder("Text will be written immediately")
		} else {
			input.SetPlaceholder("")
		}
	})
	entry.SetBorder(true)
	entry.AddOption("Immediate", nil)
	entry.AddOption("No terminator", nil)
	entry.AddOption("LF", nil)
	entry.AddOption("CR", nil)
	entry.AddOption("CRLF", nil)
	entry.AddOption("\\0", nil)
	entry.SetCurrentOption(0)
	term.widgets = append(term.widgets, entry)
	headerConfigItemRow1.AddItem(entry, 0, 1, false)

	// Add button to update config
	update := tview.NewButton("Update")
	update.SetBlurFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	update.SetSelectedFunc(term.updatePortConfig)

	headerBox.AddItem(buttonSizeFix(update, 4), 0, 1, false)
	term.widgets = append(term.widgets, update)

	// Make box to contain the output
	outputBox := tview.NewTextView()
	outputBox.SetBorder(true)
	outputBox.SetTitle("Output")
	outputBox.SetDynamicColors(true)
	term.output = outputBox
	main.AddItem(outputBox, 0, 3, false)

	// Add input field at the bottom
	input = tview.NewInputField()
	input.SetBorder(true)
	input.SetLabel("Input")
	input.SetPlaceholder("Text will be written immediately")
	input.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
		if key != tcell.KeyEnter {
			return
		}

		_, entryStr := entry.GetCurrentOption()
		terminator := ""
		switch entryStr {
		case "Immediate":
			return
		case "LF":
			terminator = "\n"
		case "CRLF":
			terminator = "\r\n"
		case "CR":
			terminator = "\r"
		case "\\0":
			terminator = "\x00"
		}
		term.Write(fmt.Sprintf("%s%s", input.GetText(), terminator), colorTx)
	})
	input.SetChangedFunc(func(str string) {
		_, entryStr := entry.GetCurrentOption()
		if len(str) > 0 && entryStr == "Immediate" {
			term.WriteLn(str, colorTx)
			input.SetText("")
		}
	})
	term.widgets = append(term.widgets, input)
	main.AddItem(input, 3, 0, false)

	// Enable mouse
	term.app.EnableMouse(true)

	// Focus on the first widget
	term.setWidget(0)

	return term
}

// Wrapper to prevent the button from becoming too large
func buttonSizeFix(button *tview.Button, height int) *tview.Flex {
	box := tview.NewFlex().SetDirection(tview.FlexRow)
	box.SetBorder(true)
	box.AddItem(button, height, 1, false)
	return box
}

func (term *Terminal) Run(channel *chan messages.Message) {
	ports := messages.SyncGetPorts(channel)
	if len(ports.Ports) > 0 {
		term.deviceSelect.SetOptions(ports.Ports, nil)
		term.deviceSelect.SetCurrentOption(ports.OpenIndex)
	}
	term.channel = channel
	term.app.Run()
	messages.SyncExit(channel)
}

func formatRxBuffer(buffer []byte, size int) string {
	return string(buffer[:size])
}

func (term *Terminal) StartReadTask(interval time.Duration) {
	go (func() {
		for {
			if term.connected && term.channel != nil {
				buffer := []byte{}
				size := 10 // Not respected at this time
				resp := messages.SyncRead(term.channel, buffer, size)
				if resp.Error != nil {
					term.WriteLn(resp.Error.Error(), colorError)
				} else {
					term.Write(formatRxBuffer(resp.Buffer, resp.Size), colorRx)
				}
			}
			term.WriteLn("Not connected...", colorDebug)
			time.Sleep(interval)
		}
	})()
}

func (term *Terminal) Write(str, color string) {
	if color != "" {
		str = color + str
	}
	term.output.Write([]byte(str))
}

func (term *Terminal) WriteLn(str, color string) {
	term.Write(fmt.Sprintf("%s\n", str), color)
}

func (term *Terminal) findSelectedWidget() int {
	for i := 0; i < len(term.widgets); i++ {
		if term.widgets[i].HasFocus() {
			return i
		}
	}
	return -1
}

func (term *Terminal) prevWidget() {
	selected := term.findSelectedWidget()
	result := 0
	if selected == -1 {
		result = 0
	} else if selected >= 1 {
		result = selected - 1
	} else {
		result = len(term.widgets) - 1
	}
	term.setWidget(result)
}

func (term *Terminal) nextWidget() {
	selected := term.findSelectedWidget()
	result := 0
	if (selected >= 0) && (selected < len(term.widgets)-1) {
		result = selected + 1
	} else {
		result = 0
	}
	term.setWidget(result)
}

func (term *Terminal) handleTab(key tcell.Key) bool {
	switch key {
	case tcell.KeyTab:
		term.nextWidget()
		return true
	case tcell.KeyBacktab:
		term.prevWidget()
		return true
	default:
		return false
	}
}

func (term *Terminal) setWidget(i int) {
	term.app.SetFocus(term.widgets[i])
}

func (term *Terminal) getPortConfig() (*messages.PortConfig, error) {
	baudInt, err := strconv.Atoi(term.baudSelect.GetText())
	if err != nil {
		return nil, err
	}
	_, dataBitsStr := term.dataBitsSelect.GetCurrentOption()
	dataBitsInt, err := strconv.Atoi(dataBitsStr)
	if err != nil {
		return nil, err
	}
	_, stopBitsStr := term.stopBitsSelect.GetCurrentOption()
	_, parityStr := term.paritySelect.GetCurrentOption()
	return &messages.PortConfig{
		BaudRate: baudInt,
		DataBits: dataBitsInt,
		StopBits: stopBitsStr,
		Parity:   parityStr,
	}, nil
}

func (term *Terminal) updatePortConfig() {
	config, err := term.getPortConfig()
	if err != nil {
		term.WriteLn(err.Error(), colorError)
	}
	_, wantedOpen := term.deviceSelect.GetCurrentOption()
	currentOpen := messages.SyncGetPorts(term.channel).OpenName
	if currentOpen == wantedOpen {
		receivedConfig := messages.SyncReconfigurePort(term.channel, config).Config
		term.WriteLn(fmt.Sprintf("Reconfigured port: %+v", *receivedConfig), colorDebug)
	} else {
		receivedConnection := messages.SyncReconnectPort(term.channel, wantedOpen, config)
		if receivedConnection.Error != nil {
			term.connected = false
			term.WriteLn(fmt.Sprintf("Failed to connect: %+v", receivedConnection), colorError)
		} else {
			term.connected = true
			term.WriteLn(fmt.Sprintf("Connected to port: %s, %+v", receivedConnection.Port, receivedConnection.Config), colorDebug)
		}
	}
}
