package ui

import (
	"fmt"
	"io"

	"github.com/gdamore/tcell/v2"
	hw "github.com/jonathangjertsen/serious/hw"
	"github.com/rivo/tview"
)

type Terminal struct {
	app          *tview.Application
	widgets      []tview.Primitive
	output       *tview.TextView
	deviceSelect *tview.DropDown
}

func NewTerminal() *Terminal {
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
	baud.SetBorder(true)
	term.widgets = append(term.widgets, baud)
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
	headerConfigItemRow1.AddItem(dataBits, 0, 1, false)

	// Add parity config
	parity := tview.NewDropDown()
	parity.SetLabel("parity ")
	parity.SetFieldWidth(6)
	parity.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	parity.SetBorder(true)
	parity.AddOption("None", nil)
	parity.AddOption("Odd", nil)
	parity.AddOption("Even", nil)
	parity.AddOption("Always 1", nil)
	parity.AddOption("Always 0", nil)
	parity.SetCurrentOption(0)
	term.widgets = append(term.widgets, parity)
	headerConfigItemRow1.AddItem(parity, 0, 1, false)

	// Add stop bits config
	stopBits := tview.NewDropDown()
	stopBits.SetLabel("stop ")
	stopBits.SetFieldWidth(6)
	stopBits.SetDoneFunc(func(key tcell.Key) {
		term.handleTab(key)
	})
	stopBits.SetBorder(true)
	stopBits.AddOption("0", nil)
	stopBits.AddOption("1", nil)
	stopBits.AddOption("2", nil)
	stopBits.SetCurrentOption(0)
	term.widgets = append(term.widgets, stopBits)
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

	headerBox.AddItem(buttonSizeFix(update, 4), 0, 1, false)
	term.widgets = append(term.widgets, update)

	// Make box to contain the output
	outputBox := tview.NewTextView()
	outputBox.SetBorder(true)
	outputBox.SetTitle("Output")
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
		io.WriteString(term, fmt.Sprintf("%s%s", input.GetText(), terminator))
	})
	input.SetChangedFunc(func(str string) {
		_, entryStr := entry.GetCurrentOption()
		if len(str) > 0 && entryStr == "Immediate" {
			term.WriteLn(str)
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

func (term *Terminal) Run() {
	term.app.Run()
}

func (term *Terminal) HwConnected(hw hw.Hw) {
	ports := hw.GetPorts()
	if len(ports) > 0 {
		term.deviceSelect.SetOptions(ports, nil)
	}
	index, selected := hw.Selected()
	if selected != nil {
		term.deviceSelect.SetCurrentOption(index)
	} else {
		term.deviceSelect.SetCurrentOption(0)
	}
}

func (term *Terminal) Write(str []byte) (n int, err error) {
	return term.output.Write(str)
}

func (term *Terminal) WriteLn(str string) {
	io.WriteString(term, fmt.Sprintf("%s\n", str))
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
