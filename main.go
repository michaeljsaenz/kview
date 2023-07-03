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
	clientset, config := k8s.GetClientSet()

	// retrieve namespaces
	namespaceList := k8s.GetNamespaces(*clientset)

	// get current cluster context
	currentContext := k8s.GetCurrentContext()

	// create a new app, window title and size
	app := app.New()
	win := app.NewWindow("KView")
	win.SetMaster()
	win.Resize(fyne.NewSize(1200, 700))
	win.CenterOnScreen()

	// list binding, bind pod list (podData) to data
	var podData []string
	data, list := ui.GetListData(&podData)

	// intial/base widgets and windows
	topWindowLabel, topWindow, rightWindow, rightWindowTitle := ui.CreateWindows(currentContext)

	podStatus, input, listTitle := ui.CreateBaseWidgets()

	podLabelsLabel, podLabels, podLabelsScroll, podAnnotationsLabel, podAnnotations, podAnnotationsScroll,
		podEventsLabel, podEvents, podEventsScroll, podLogsLabel, podLog, podLogScroll, podDetailLabel, podDetailLog, podDetailScroll,
		podVolumesLabel, podVolumes, podVolumesScroll := ui.CreateBaseTabs()

	podTabs, podLogTabs := ui.CreateBaseTabContainers(podLabelsLabel, podLabelsScroll, podAnnotationsLabel, podAnnotationsScroll,
		podEventsLabel, podEventsScroll, podLogsLabel, podLogScroll, podDetailLabel, podDetailScroll, podVolumesLabel, podVolumesScroll)

	// create the namespace dropdown list widget
	namespaceListDropdown := widget.NewSelect(namespaceList, func(selectedNamespace string) {
		if selectedNamespace != "" {
			podData = k8s.GetPodDataWithNamespace(*clientset, selectedNamespace)
		}
		input.Text = ""
		input.Refresh()
		data.Reload()
		list.UnselectAll()

	})

	namespaceListDropdown.PlaceHolder = "Select namespace..."
	namespaceListDropdown.FocusGained()

	// refresh and clear pod list data
	refresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		if namespaceListDropdown.Selected == "" {
			podData = []string{}
		} else {
			podData = k8s.GetPodDataWithNamespace(*clientset, namespaceListDropdown.Selected)
		}
		input.Text = ""
		input.Refresh()
		data.Reload()
		list.UnselectAll()
		podTabs.SelectIndex(0)
		podLogTabs.SelectIndex(0)
		// remove container log tabs before loading current selection
		podLogTabItems := len(podLogTabs.Items)
		for podLogTabItems > 1 {
			for _, item := range podLogTabs.Items {
				if item.Text != podLogsLabel.Text {
					podLogTabs.Remove(item)
				}
			}
			podLogTabItems = len(podLogTabs.Items)
		}
		podStatus.Text = "Status: \n" + "Age: \n" + "Namespace: \n" + "Node: "
		podStatus.Refresh()
		rightWindowTitle.Text = "Select application (pod)..."
		rightWindowTitle.Refresh()
	})

	stringErrorResponse, errorPresent := utils.CheckForError(podData)

	if !errorPresent {
		stringErrorResponse, errorPresent = utils.CheckForError(namespaceList)
	}

	if errorPresent {
		namespaceListDropdown.Disable()
		input.Disable()
		rightWindowTitle, list, refresh = ui.SetupErrorUI(stringErrorResponse, list)
	}

	// search application name (input list field)
	if !errorPresent {
		input.OnSubmitted = func(s string) {
			podData = ui.InputOnSubmitted(input, *clientset, namespaceListDropdown)
			data.Reload()
			list.UnselectAll()
		}
	}

	yamlButton := ui.CreateIconButton("Application (Pod) YAML", theme.ZoomInIcon())
	yamlButton.Hide()

	execButtons := ui.CreateBaseExecIconButton("", theme.LoginIcon())
	for _, execButton := range execButtons {
		execButton.Hide()
	}

	gridOne := container.New(layout.NewGridLayout(1), yamlButton)
	gridTwo := container.New(layout.NewGridLayoutWithColumns(2), execButtons[0], execButtons[1], execButtons[2], execButtons[3], execButtons[4],
		execButtons[5], execButtons[6], execButtons[7], execButtons[8], execButtons[9])

	ui.ListOnSelected(list, data, *clientset, *config, rightWindowTitle, podStatus, podLabels,
		podAnnotations, podEvents, podVolumes, podLog, podDetailLog, podTabs, podLogTabs, podLogScroll,
		podLogsLabel, app, yamlButton, execButtons, namespaceListDropdown)

	//return tabs to initial tab (index 0)
	list.OnUnselected = func(id widget.ListItemID) {
		podTabs.SelectIndex(0)
		podLogTabs.SelectIndex(0)
		for _, execButton := range execButtons {
			execButton.Hide()
		}
	}

	rightContainer := container.NewBorder(
		container.NewVBox(rightWindowTitle, podStatus, podTabs, podLogTabs, gridOne, gridTwo),
		nil, nil, nil, rightWindow)

	listContainer := container.NewBorder(container.NewVBox(listTitle, namespaceListDropdown, input),
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
