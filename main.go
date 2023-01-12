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

	// get a list of all pods
	podData := getPodData(*clientset)

	// get current cluster context
	currentContext := getCurrentContext()

	// create a new app
	app := app.New()
	// create a new window with app title
	win := app.NewWindow("KUI")
	// resize fyne app window
	win.Resize(fyne.NewSize(1200, 700)) // first width, then height

	// list binding, bind pod list data(podData) to data
	data, list := getListData(&podData)

	// top window label
	topLabel := canvas.NewText(("Cluster Context: " + currentContext), color.NRGBA{R: 57, G: 112, B: 228, A: 255})
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

	// setup pod tabs
	podTabs := container.NewAppTabs(
		container.NewTabItem(podLabelsLabel.Text, podLabelsScroll),
		container.NewTabItem(podAnnotationsLabel.Text, podAnnotationsScroll),
		container.NewTabItem(podEventsLabel.Text, podEventsScroll),
	)

	// setup pod log tabs
	podLogsLabel := widget.NewLabel("")
	podLogsLabel.TextStyle = fyne.TextStyle{Monospace: true}
	defaultTabItem := container.NewTabItem("Logs", podLogsLabel)
	podLogTabs := container.NewAppTabs(defaultTabItem)

	// update pod list data
	refresh := widget.NewButton("Refresh", func() {
		podData = getPodData(*clientset)
		list.UnselectAll()
		data.Reload()
	})

	list.OnSelected = func(id widget.ListItemID) {
		selectedPod, err := data.GetValue(id)
		if err != nil {
			panic(err.Error())
		}
		title.Text = "Application (Pod): " + selectedPod
		title.Refresh()

		newPodStatus, newPodAge, newPodNamespace, newPodLabels, newPodAnnotations, newNodeName, newContainers := getPodDetail(*clientset, selectedPod)

		podStatus.Text = "Status: " + newPodStatus + "\n" +
			"Age: " + newPodAge + "\n" +
			"Namespace: " + newPodNamespace + "\n" +
			"Node: " + newNodeName
		podStatus.Refresh()

		podLabels.Text = newPodLabels
		podLabels.Refresh()

		podAnnotations.Text = newPodAnnotations
		podAnnotations.Refresh()

		// get pod events
		newPodEvents := getPodEvents(*clientset, selectedPod)
		strNewPodEvents := strings.Join(newPodEvents, "\n")
		podEvents.Text = strNewPodEvents
		podEvents.Refresh()

		for _, tabContainerName := range newContainers {
			podLogStream := getPodLogs(*clientset, newPodNamespace, selectedPod, tabContainerName)
			podLog := widget.NewLabel(podLogStream)
			podLog.TextStyle = fyne.TextStyle{Monospace: true}
			podLog.Wrapping = fyne.TextWrapBreak
			podLogScroll := container.NewScroll(podLog)
			podLogScroll.SetMinSize(fyne.Size{Height: 200})
			podLogTabs.Append(container.NewTabItem(tabContainerName, podLogScroll))
			podLog.Refresh()
		}
		podLogTabs.Refresh()
	}

	list.OnUnselected = func(id widget.ListItemID) {
		for _, tabItem := range podLogTabs.Items {
			if tabItem != defaultTabItem {
				podLogTabs.Remove(tabItem)
			}
		}
	}

	rightContainer := container.NewBorder(
		container.NewVBox(title, podStatus, podTabs, podLogTabs),
		nil, nil, nil, rightWinContent)

	listTitle := widget.NewLabel("Application (Pod)")
	listTitle.Alignment = fyne.TextAlignCenter
	listTitle.TextStyle = fyne.TextStyle{Monospace: true}

	// search application name (input list field)
	input := widget.NewEntry()
	input.SetPlaceHolder("Search application...")
	// submit to func input string (pod name), return new pod list
	input.OnSubmitted = func(s string) {
		inputText := input.Text
		var inputTextList []string
		if inputText == "" {
			podData = getPodData(*clientset)
			data.Reload()
			list.UnselectAll()
		} else {
			for _, pod := range podData {
				if strings.Contains(pod, inputText) {
					inputTextList = append(inputTextList, pod)
				}
			}
			podData = inputTextList
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

//TODO catch panic when cluster context not available:
// panic: Get "https://1.2.3.4:443/api/v1/pods": dial tcp 1.2.3.4:443: i/o timeout

//TODO test if kubeConfig not accessible/ not set
//TODO test if clusterContext not set / empty
//TODO add copy capability
//TODO clear podTab data on refresh, similar to podLogTab data on refresh
//TODO optimize the log tabs
