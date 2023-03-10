package main

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// setup k8s clientset
	clientset := getClientSet()

	// load initial pod list data (UI `list`)
	podData := getPodData(*clientset)

	// get current cluster context
	currentContext := getCurrentContext()

	// create a new app
	app := app.New()
	// create a new window with app title
	win := app.NewWindow("KView")
	// resize fyne app window
	win.Resize(fyne.NewSize(1200, 700)) // first width, then height

	// list binding, bind pod list (podData) to data
	data, list := getListData(&podData)

	// top window label
	topLabel := canvas.NewText(("Cluster Context: " + currentContext),
		color.NRGBA{R: 57, G: 112, B: 228, A: 255})
	topLabel.TextStyle = fyne.TextStyle{Monospace: true}
	topContent := container.New(layout.NewCenterLayout(), topLabel)

	// right side of split
	rightWinContent := container.NewMax()
	title := widget.NewLabel("Select application (pod)...")
	title.TextStyle = fyne.TextStyle{Monospace: true}

	// pod status
	podStatus := widget.NewLabel("")
	podStatus.TextStyle = fyne.TextStyle{Monospace: true}

	// get pod labels, annotations, events for tabs
	podLabelsLabel, podLabels, podLabelsScroll := getPodTabData("Labels")
	podAnnotationsLabel, podAnnotations, podAnnotationsScroll := getPodTabData("Annotations")
	podEventsLabel, podEvents, podEventsScroll := getPodTabData("Events")
	podLogsLabel, podLog, podLogScroll := getPodTabData("")

	// setup tabs
	podTabs := container.NewAppTabs(
		container.NewTabItem(podLabelsLabel.Text, podLabelsScroll),
		container.NewTabItem(podAnnotationsLabel.Text, podAnnotationsScroll),
		container.NewTabItem(podEventsLabel.Text, podEventsScroll),
	)
	podLogTabs := container.NewAppTabs(
		container.NewTabItem(podLogsLabel.Text, podLogScroll),
	)

	// setup input widget
	input := widget.NewEntry()
	input.SetPlaceHolder("Search application (pod)...")

	// update pod list data
	refresh := widget.NewButton("Refresh", func() {
		podData = getPodData(*clientset)
		list.UnselectAll()
		data.Reload()
		input.Text = ""
		input.Refresh()
	})

	// check for error from initial pod list data
	stringErrorResponse, errorPresent := checkForError(podData)
	if errorPresent {
		title, list, refresh = setupErrorUI(stringErrorResponse, list)
	}
	listOnSelected(list, data, *clientset, title, podStatus, podLabels,
		podAnnotations, podEvents, podLog, podLogTabs, podLogScroll)

	rightContainer := container.NewBorder(
		container.NewVBox(title, podStatus, podTabs, podLogTabs),
		nil, nil, nil, rightWinContent)

	listTitle := widget.NewLabel("Application (Pod)")
	listTitle.Alignment = fyne.TextAlignCenter
	listTitle.TextStyle = fyne.TextStyle{Monospace: true}

	// search application name (input list field)
	if !errorPresent {
		input.OnSubmitted = func(s string) {
			podData = inputOnSubmitted(input, *clientset, podData)
			data.Reload()
			list.UnselectAll()
		}
	}

	listContainer := container.NewBorder(container.NewVBox(listTitle, input), nil, nil, nil, list)

	// podData(list) left side, podData detail right side
	split := container.NewHSplit(listContainer, rightContainer)
	split.Offset = 0.3

	// check current cluster context to update top window label
	go func() {
		for range time.Tick(time.Second * 5) {
			currentContext = getCurrentContext()
			if strings.Contains(topLabel.Text, currentContext) {
				continue
			} else {
				topLabel.Text = ("Cluster Context: " + currentContext)
				topLabel.Refresh()
			}
		}
	}()

	win.SetContent(container.NewBorder(topContent, refresh, nil, nil, split))
	win.ShowAndRun()
}
