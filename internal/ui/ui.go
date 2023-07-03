package ui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/michaeljsaenz/kview/internal/k8s"
	"github.com/michaeljsaenz/kview/internal/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	refresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {})
	refresh.Refresh()
	return title, list, refresh
}

func ListOnSelected(list *widget.List, data binding.ExternalStringList, clientset kubernetes.Clientset, config rest.Config, title, podStatus,
	podLabels, podAnnotations, podEvents, podVolumes, podLog *widget.Label, podDetailLog *widget.Label, podTabs *container.AppTabs, podLogTabs *container.AppTabs,
	podLogScroll *container.Scroll, podLogsLabel *widget.Label, app fyne.App, yb *widget.Button, execButtons []*widget.Button, namespaceListDropdown *widget.Select) {
	list.OnSelected = func(id widget.ListItemID) {

		selectedPod, err := data.GetValue(id)
		if err != nil {
			panic(err.Error())
		}
		title.Text = "Application (Pod): " + selectedPod
		title.Refresh()

		newPodStatus, newPodAge, newNodeName, newContainers := k8s.GetPodDetail(clientset, selectedPod, namespaceListDropdown.Selected)

		newPodNamespace := namespaceListDropdown.Selected
		podStatus.Text = "Status: " + newPodStatus + "\n" +
			"Age: " + newPodAge + "\n" +
			"Namespace: " + newPodNamespace + "\n" +
			"Node: " + newNodeName
		podStatus.Refresh()

		podTabs.OnSelected = func(tabItemName *container.TabItem) {
			selectedTabItemName := tabItemName.Text
			switch selectedTabItemName {
			case "Labels":
				// get pod labels
				newPodLabels := k8s.GetPodLabels(clientset, selectedPod, namespaceListDropdown.Selected)
				podLabels.Text = newPodLabels
				podLabels.Refresh()
			case "Annotations":
				// get pod annotations
				newPodAnnotations := k8s.GetPodAnnotations(clientset, selectedPod, namespaceListDropdown.Selected)
				podAnnotations.Text = newPodAnnotations
				podAnnotations.Refresh()
			case "Events":
				// get pod events
				newPodEvents := k8s.GetPodEvents(clientset, selectedPod, namespaceListDropdown.Selected)
				strNewPodEvents := strings.Join(newPodEvents, "\n")
				podEvents.Text = strNewPodEvents
				podEvents.Refresh()
			case "Volumes":
				// get pod volumes
				newVolumes, err := k8s.GetPodVolumes(clientset, selectedPod, namespaceListDropdown.Selected)
				if err != nil {
					fmt.Printf("error with GetPodVolumes: %v\n", err)
				}
				podVolumes.Text = newVolumes
				podVolumes.Refresh()
			}

		}

		currentSelectedTabItemName := podTabs.Selected()
		switch currentSelectedTabItemName.Text {
		case "Labels":
			// get pod labels
			newPodLabels := k8s.GetPodLabels(clientset, selectedPod, namespaceListDropdown.Selected)
			podLabels.Text = newPodLabels
			podLabels.Refresh()
		case "Annotations":
			// get pod annotations
			newPodAnnotations := k8s.GetPodAnnotations(clientset, selectedPod, namespaceListDropdown.Selected)
			podAnnotations.Text = newPodAnnotations
			podAnnotations.Refresh()
		case "Events":
			// get pod events
			newPodEvents := k8s.GetPodEvents(clientset, selectedPod, namespaceListDropdown.Selected)
			strNewPodEvents := strings.Join(newPodEvents, "\n")
			podEvents.Text = strNewPodEvents
			podEvents.Refresh()
		case "Volumes":
			// get pod volumes
			newVolumes, err := k8s.GetPodVolumes(clientset, selectedPod, namespaceListDropdown.Selected)
			if err != nil {
				fmt.Printf("error with GetPodVolumes: %v\n", err)
			}
			podVolumes.Text = newVolumes
			podVolumes.Refresh()
		}

		yb.Show()

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

		for _, tabContainerName := range newContainers {
			podLogScroll.SetMinSize(fyne.Size{Height: 200})
			podLogTabs.Append(container.NewTabItemWithIcon(tabContainerName, theme.DocumentIcon(), podLogScroll))
			podLogTabs.Refresh()
		}

		for i, buttonContainerName := range newContainers {
			execButtons[i].SetText(buttonContainerName)
			execButtons[i].Show()
			if i == 9 {
				break
			}
		}

		// assign the OnTapped function to each button
		for _, button := range execButtons {
			button := button
			button.OnTapped = func() {

				win := app.NewWindow("Container Name: " + button.Text)

				// Create an entry field for user input
				entry := widget.NewEntry()
				entry.SetPlaceHolder("Enter a command")

				// Create a label for displaying the prompt and command output
				outputLabel := widget.NewLabel("")
				podOutputScroll := container.NewScroll(outputLabel)

				// OnSubmitted event handler to execute the command
				entry.OnSubmitted = func(command string) {
					client := k8s.GetClientInterface(clientset)
					// execute the command and return string output
					commandOutput, err := k8s.ExecCmd(client, config, selectedPod, button.Text, namespaceListDropdown.Selected, command, nil)
					if err != nil {
						fmt.Println(err)
					}

					// update the output label
					commandOutput = utils.RemoveANSIEscapeCodes(commandOutput)
					outputLabel.SetText(commandOutput)

					// Use the captured stdout content as needed
					podOutputScroll = container.NewScroll(outputLabel)
					entry.SetText("") // clear the input field
				}

				bottomBox := container.NewVBox(
					widget.NewButtonWithIcon("Copy Output", theme.ContentCopyIcon(), func() {
						win.Clipboard().SetContent(outputLabel.Text)
					}),
				)

				content := container.NewBorder(entry, bottomBox, nil, nil, podOutputScroll)

				win.SetContent(content)
				win.Resize(fyne.NewSize(1200, 700))
				win.Show()
			}
		}

		podLogTabs.OnSelected = func(containerTabItemName *container.TabItem) {
			var containerLogStream string
			if podLogsLabel.Text != podLogTabs.Selected().Text {
				containerLogStream = k8s.GetPodLogs(clientset, newPodNamespace, selectedPod, containerTabItemName.Text)
			}
			podLog = widget.NewLabel(containerLogStream)
			podLogScroll = container.NewScroll(podLog)
			podLogScroll.SetMinSize(fyne.Size{Height: 200})
			podLogTabs.Selected().Content = podLogScroll
			podLogScroll.Refresh()
			podLogTabs.Refresh()
		}

		yb.OnTapped = func() {
			// export yaml and display in new window
			win := app.NewWindow("Application (Pod): " + selectedPod)
			podYaml, err := k8s.GetPodYaml(clientset, newPodNamespace, selectedPod)
			if err != nil {
				fmt.Printf("error with GetPodYaml: %v\n", err)
			}
			podYamlScroll := container.NewScroll(widget.NewLabel(podYaml))

			bottomBox := container.NewVBox(
				widget.NewButtonWithIcon("Copy YAML", theme.ContentCopyIcon(), func() {
					win.Clipboard().SetContent(podYaml)
				}),
			)
			content := container.NewBorder(nil, bottomBox, nil, nil, podYamlScroll)

			win.SetContent(content)
			win.Resize(fyne.NewSize(1200, 700))
			win.Show()
		}

	}
}

func InputOnSubmitted(input *widget.Entry, clientset kubernetes.Clientset, namespaceListDropdown *widget.Select) []string {
	// submit to func input string (pod name), return new pod list
	inputText := input.Text
	var inputTextList, podData []string
	if inputText == "" {
		return podData
	}
	if namespaceListDropdown.Selected != "" {
		podData = k8s.GetPodDataWithNamespace(clientset, namespaceListDropdown.Selected)
		for _, pod := range podData {
			if strings.Contains(pod, inputText) {
				inputTextList = append(inputTextList, pod)
			}
		}
		podData = inputTextList
		return podData
	} else {
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
	podStatus := widget.NewLabel("Status: \n" + "Age: \n" + "Namespace: \n" + "Node: ")
	podStatus.TextStyle = fyne.TextStyle{Monospace: true}

	// setup input widget
	input := widget.NewEntry()
	input.SetPlaceHolder("Search application (pod)...")

	listTitle := widget.NewLabel("Application (Pod)")
	listTitle.Alignment = fyne.TextAlignCenter
	listTitle.TextStyle = fyne.TextStyle{Monospace: true}

	return podStatus, input, listTitle
}

func CreateBaseTabs() (*widget.Label, *widget.Label, *container.Scroll, *widget.Label, *widget.Label, *container.Scroll,
	*widget.Label, *widget.Label, *container.Scroll, *widget.Label, *widget.Label, *container.Scroll, *widget.Label, *widget.Label,
	*container.Scroll, *widget.Label, *widget.Label, *container.Scroll) {

	//get pod labels, annotations, events, volumes for tabs
	podDetailLabel, podDetailLog, podDetailScroll := GetPodTabData("")
	podLabelsLabel, podLabels, podLabelsScroll := GetPodTabData("Labels")
	podAnnotationsLabel, podAnnotations, podAnnotationsScroll := GetPodTabData("Annotations")
	podEventsLabel, podEvents, podEventsScroll := GetPodTabData("Events")
	podVolumesLabel, podVolumes, podVolumesScroll := GetPodTabData("Volumes")
	podLogsLabel, podLog, podLogScroll := GetPodTabData("")

	return podLabelsLabel, podLabels, podLabelsScroll, podAnnotationsLabel, podAnnotations, podAnnotationsScroll,
		podEventsLabel, podEvents, podEventsScroll, podLogsLabel, podLog, podLogScroll, podDetailLabel, podDetailLog, podDetailScroll,
		podVolumesLabel, podVolumes, podVolumesScroll

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
	podEventsLabel *widget.Label, podEventsScroll *container.Scroll, podLogsLabel *widget.Label, podLogScroll *container.Scroll, podDetailLabel *widget.Label, podDetailScroll *container.Scroll,
	podVolumesLabel *widget.Label, podVolumesScroll *container.Scroll) (*container.AppTabs, *container.AppTabs) {
	podTabs := container.NewAppTabs(
		container.NewTabItemWithIcon(podDetailLabel.Text, theme.MailForwardIcon(), podDetailScroll),
		container.NewTabItem(podLabelsLabel.Text, podLabelsScroll),
		container.NewTabItem(podAnnotationsLabel.Text, podAnnotationsScroll),
		container.NewTabItem(podEventsLabel.Text, podEventsScroll),
		container.NewTabItem(podVolumesLabel.Text, podVolumesScroll),
	)
	podLogTabs := container.NewAppTabs(
		container.NewTabItemWithIcon(podLogsLabel.Text, theme.MailForwardIcon(), podLogScroll),
	)

	return podTabs, podLogTabs
}

func CreateIconButton(buttonName string, iconTheme fyne.Resource) *widget.Button {
	button := widget.NewButtonWithIcon(buttonName, iconTheme, func() {})
	return button
}

func CreateBaseExecIconButton(buttonName string, iconTheme fyne.Resource) []*widget.Button {
	buttons := make([]*widget.Button, 10)

	for i := 0; i < 10; i++ {
		button := widget.NewButtonWithIcon(buttonName, iconTheme, func() {})
		buttons[i] = button
	}

	return buttons
}
