package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/michaeljsaenz/kview/internal/k8s"
	"github.com/michaeljsaenz/kview/internal/ui"
	"github.com/michaeljsaenz/kview/internal/utils"
)

func main() {
	// setup k8s clientset
	clientset := k8s.GetClientSet()

	// load initial pod list data (UI `list`)
	podData := k8s.GetPodData(*clientset)

	// get current cluster context
	currentContext := k8s.GetCurrentContext()

	// create a new app, window title and size
	app := app.New()
	win := app.NewWindow("KView")
	// icon, _ := fyne.LoadResourceFromPath("internal/assets/icon/icon.png")
	// win.SetIcon(icon)
	win.SetMaster()
	win.Resize(fyne.NewSize(1200, 700))
	win.CenterOnScreen()

	// list binding, bind pod list (podData) to data
	data, list := ui.GetListData(&podData)

	// intial/base widgets and windows
	topWindowLabel, topWindow, rightWindow, rightWindowTitle := ui.CreateWindows(currentContext)

	podStatus, input, listTitle := ui.CreateBaseWidgets()

	podLabelsLabel, podLabels, podLabelsScroll, podAnnotationsLabel, podAnnotations, podAnnotationsScroll,
		podEventsLabel, podEvents, podEventsScroll, podLogsLabel, podLog, podLogScroll := ui.CreateBaseTabs()

	podTabs, podLogTabs := ui.CreateBaseTabContainers(podLabelsLabel, podLabelsScroll, podAnnotationsLabel, podAnnotationsScroll,
		podEventsLabel, podEventsScroll, podLogsLabel, podLogScroll)

	// update pod list data
	refresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		podData = ui.RefreshButton(input, *clientset)
		data.Reload()
		list.UnselectAll()
	})

	stringErrorResponse, errorPresent := utils.CheckForError(podData)
	if errorPresent {
		rightWindowTitle, list, refresh = ui.SetupErrorUI(stringErrorResponse, list)
	}

	// search application name (input list field)
	if !errorPresent {
		input.OnSubmitted = func(s string) {
			podData = ui.InputOnSubmitted(input, *clientset)
			data.Reload()
			list.UnselectAll()
		}
	}
	describeButton := widget.NewButtonWithIcon("Container Detail", theme.ZoomInIcon(), func() {
	})
	describeButton.Hide()

	yamlButton := widget.NewButtonWithIcon("YAML", theme.ZoomInIcon(), func() {
	})
	yamlButton.Hide()

	grid := container.New(layout.NewGridLayout(1), describeButton, yamlButton)

	ui.ListOnSelected(list, data, *clientset, rightWindowTitle, podStatus, podLabels,
		podAnnotations, podEvents, podLog, podLogTabs, podLogScroll, app, describeButton, yamlButton)

	rightContainer := container.NewBorder(
		container.NewVBox(rightWindowTitle, podStatus, podTabs, podLogTabs, grid),
		nil, nil, nil, rightWindow)

	listContainer := container.NewBorder(container.NewVBox(listTitle, input),
		nil, nil, nil, list)

	// podData(list) left side, podData detail right side
	split := container.NewHSplit(listContainer, rightContainer)
	split.Offset = 0.3

	// check current cluster context to update top window label
	go func() {
		for range time.Tick(time.Second * 5) {
			currentContext = k8s.GetCurrentContext()
			if strings.Contains(topWindowLabel.Text, currentContext) {
				continue
			} else {
				topWindowLabel.Text = ("Cluster Context: " + currentContext)
				topWindowLabel.Refresh()
			}
		}
	}()

	win.SetContent(container.NewBorder(topWindow, refresh, nil, nil, split))
	win.ShowAndRun()
}
