package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"k8s.io/client-go/kubernetes"
)

func getListData(podData *[]string) (binding.ExternalStringList, *widget.List) {
	// list binding, bind pod list data to data
	data := binding.BindStringList(
		podData,
	)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	return data, list
}

func setupErrorUI(stringErrorResponse string, list *widget.List) (*widget.Label, *widget.List, *widget.Button) {
	title := widget.NewLabel("")
	title.TextStyle = fyne.TextStyle{Monospace: true}
	title.Text = stringErrorResponse
	title.Wrapping = fyne.TextWrapBreak
	title.TextStyle = fyne.TextStyle{Italic: true, Bold: true}
	title.Refresh()
	list.Hide()
	refresh := widget.NewButton("Refresh", func() {})
	refresh.Refresh()
	return title, list, refresh
}

func listOnSelected(list *widget.List, data binding.ExternalStringList, clientset kubernetes.Clientset, title, podStatus,
	podLabels, podAnnotations, podEvents, podLog *widget.Label, podLogTabs *container.AppTabs, podLogScroll *container.Scroll) {
	list.OnSelected = func(id widget.ListItemID) {
		selectedPod, err := data.GetValue(id)
		if err != nil {
			panic(err.Error())
		}
		title.Text = "Application (Pod): " + selectedPod
		title.Refresh()

		podNamespace := getPodNamespace(clientset, selectedPod)
		newPodStatus, newPodAge, newPodNamespace, newPodLabels, newPodAnnotations, newNodeName, newContainers := getPodDetail(clientset, selectedPod, podNamespace)

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
		newPodEvents := getPodEvents(clientset, selectedPod, podNamespace)
		strNewPodEvents := strings.Join(newPodEvents, "\n")
		podEvents.Text = strNewPodEvents
		podEvents.Refresh()

		// remove container log tabs before loading current selection
		podTabItems := len(podLogTabs.Items)
		for podTabItems > 0 {
			for _, item := range podLogTabs.Items {
				podLogTabs.Remove(item)
			}
			podTabItems = len(podLogTabs.Items)
		}

		for _, tabContainerName := range newContainers {
			podLogStream := getPodLogs(clientset, newPodNamespace, selectedPod, tabContainerName)
			podLog = widget.NewLabel(podLogStream)
			podLogScroll = container.NewScroll(podLog)
			podLogScroll.SetMinSize(fyne.Size{Height: 200})
			podLogTabs.Append(container.NewTabItem(tabContainerName, podLogScroll))
			podLogTabs.Refresh()
		}
	}
}

func inputOnSubmitted(input *widget.Entry,
	clientset kubernetes.Clientset, podData []string) []string {
	// submit to func input string (pod name), return new pod list
	inputText := input.Text
	var inputTextList []string
	if inputText == "" {
		podData = getPodData(clientset)
		return podData
	} else {
		for _, pod := range podData {
			if strings.Contains(pod, inputText) {
				inputTextList = append(inputTextList, pod)
			}
		}
		podData = inputTextList
		return podData
	}
}
