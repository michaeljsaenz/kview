package ui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/michaeljsaenz/kview/internal/k8s"
	"k8s.io/client-go/kubernetes"
)

func GetListData(podData *[]string) (binding.ExternalStringList, *widget.List) {
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

func SetupErrorUI(stringErrorResponse string, list *widget.List) (*widget.Label, *widget.List, *widget.Button) {
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

func ListOnSelected(list *widget.List, data binding.ExternalStringList, clientset kubernetes.Clientset, title, podStatus,
	podLabels, podAnnotations, podEvents, podLog *widget.Label, podLogTabs *container.AppTabs, podLogScroll *container.Scroll) {
	list.OnSelected = func(id widget.ListItemID) {
		selectedPod, err := data.GetValue(id)
		if err != nil {
			panic(err.Error())
		}
		title.Text = "Application (Pod): " + selectedPod
		title.Refresh()

		podNamespace := k8s.GetPodNamespace(clientset, selectedPod)
		newPodStatus, newPodAge, newPodNamespace, newPodLabels, newPodAnnotations, newNodeName, newContainers := k8s.GetPodDetail(clientset, selectedPod, podNamespace)

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
		newPodEvents := k8s.GetPodEvents(clientset, selectedPod, podNamespace)
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
			podLogStream := k8s.GetPodLogs(clientset, newPodNamespace, selectedPod, tabContainerName)
			podLog = widget.NewLabel(podLogStream)
			podLogScroll = container.NewScroll(podLog)
			podLogScroll.SetMinSize(fyne.Size{Height: 200})
			podLogTabs.Append(container.NewTabItem(tabContainerName, podLogScroll))
			podLogTabs.Refresh()
		}
	}
}

func InputOnSubmitted(input *widget.Entry, clientset kubernetes.Clientset) []string {
	// submit to func input string (pod name), return new pod list
	inputText := input.Text
	var inputTextList []string
	podData := k8s.GetPodData(clientset)
	if inputText == "" {
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

func CreateWindows(currentContext string) (*canvas.Text, *fyne.Container, *fyne.Container, *widget.Label) {

	// top window label
	topWindowLabel := canvas.NewText(("Cluster Context: " + currentContext), color.NRGBA{R: 57, G: 112, B: 228, A: 255})
	topWindowLabel.TextStyle = fyne.TextStyle{Monospace: true}
	topWindow := container.New(layout.NewCenterLayout(), topWindowLabel)

	// right side of split
	rightWindow := container.NewMax()
	title := widget.NewLabel("Select application (pod)...")
	title.TextStyle = fyne.TextStyle{Monospace: true}

	return topWindowLabel, topWindow, rightWindow, title

}

func CreateBaseWidgets() (*widget.Label, *widget.Entry, *widget.Label) {
	// setup pod status
	podStatus := widget.NewLabel("")
	podStatus.TextStyle = fyne.TextStyle{Monospace: true}

	// setup input widget
	input := widget.NewEntry()
	input.SetPlaceHolder("Search application (pod)...")

	listTitle := widget.NewLabel("Application (Pod)")
	listTitle.Alignment = fyne.TextAlignCenter
	listTitle.TextStyle = fyne.TextStyle{Monospace: true}

	return podStatus, input, listTitle
}
func RefreshButton(input *widget.Entry, clientset kubernetes.Clientset) []string {
	podData := k8s.GetPodData(clientset)
	input.Text = ""
	input.Refresh()
	return podData
}

func CreateBaseTabs() (*widget.Label, *widget.Label, *container.Scroll, *widget.Label, *widget.Label, *container.Scroll,
	*widget.Label, *widget.Label, *container.Scroll, *widget.Label, *widget.Label, *container.Scroll) {

	//get pod labels, annotations, events for tabs
	podLabelsLabel, podLabels, podLabelsScroll := GetPodTabData("Labels")
	podAnnotationsLabel, podAnnotations, podAnnotationsScroll := GetPodTabData("Annotations")
	podEventsLabel, podEvents, podEventsScroll := GetPodTabData("Events")
	podLogsLabel, podLog, podLogScroll := GetPodTabData("")

	return podLabelsLabel, podLabels, podLabelsScroll, podAnnotationsLabel, podAnnotations, podAnnotationsScroll,
		podEventsLabel, podEvents, podEventsScroll, podLogsLabel, podLog, podLogScroll

}

func GetPodTabData(widgetLabelName string) (widgetNameLabel *widget.Label, widgetName *widget.Label, widgetNameScroll *container.Scroll) {
	widgetNameLabel = widget.NewLabel(widgetLabelName)
	widgetNameLabel.TextStyle = fyne.TextStyle{Monospace: true}
	widgetName = widget.NewLabel("")
	widgetName.TextStyle = fyne.TextStyle{Monospace: true}
	widgetNameScroll = container.NewScroll(widgetName)
	widgetNameScroll.SetMinSize(fyne.Size{Height: 100})
	return widgetNameLabel, widgetName, widgetNameScroll
}

func CreateBaseTabContainers(podLabelsLabel *widget.Label, podLabelsScroll *container.Scroll, podAnnotationsLabel *widget.Label, podAnnotationsScroll *container.Scroll,
	podEventsLabel *widget.Label, podEventsScroll *container.Scroll, podLogsLabel *widget.Label, podLogScroll *container.Scroll) (*container.AppTabs, *container.AppTabs) {
	podTabs := container.NewAppTabs(
		container.NewTabItem(podLabelsLabel.Text, podLabelsScroll),
		container.NewTabItem(podAnnotationsLabel.Text, podAnnotationsScroll),
		container.NewTabItem(podEventsLabel.Text, podEventsScroll),
	)
	podLogTabs := container.NewAppTabs(
		container.NewTabItem(podLogsLabel.Text, podLogScroll),
	)

	return podTabs, podLogTabs
}
