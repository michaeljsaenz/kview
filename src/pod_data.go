package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// get pod names to populate initial list
func getPodData(c kubernetes.Clientset) (podData []string) {
	pods, err := c.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, pod := range pods.Items {
		podData = append(podData, pod.Name)
	}
	return podData
}

func getPodDetail(c kubernetes.Clientset, listItemID int, selectedPod string) (string, string, string, string, string, string, []string) {
	pods, err := c.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		if pod.Name == selectedPod {
			var containers []string
			for _, container := range pod.Spec.Containers {
				containers = append(containers, container.Name)
			}
			podCreationTime := pod.GetCreationTimestamp()
			age := time.Since(podCreationTime.Time).Round(time.Second)

			return string(pod.Status.Phase), age.String(), string(pod.Namespace), convertMapToString(pod.Labels), convertMapToString(pod.Annotations),
				pod.Spec.NodeName, containers
		}
	}
	return "", "", "", "", "", "", []string{}
}

func getPodEvents(c kubernetes.Clientset, selectedPod string) (podEvents []string) {
	events, _ := c.CoreV1().Events("").List(context.TODO(), v1.ListOptions{FieldSelector: fmt.Sprintf("involvedObject.name=%s", selectedPod), TypeMeta: v1.TypeMeta{Kind: "Pod"}})
	for _, item := range events.Items {
		podEvents = append(podEvents, "~> "+item.EventTime.Time.Format("2006-01-02 15:04:05")+", "+item.Message)
	}
	return podEvents
}

func getPodLogs(c kubernetes.Clientset, podNamespace string, selectedPod string, containerName string) (podLog string) {
	podLogReq := c.CoreV1().Pods(podNamespace).GetLogs(selectedPod, &corev1.PodLogOptions{Container: containerName})
	podStream, err := podLogReq.Stream(context.TODO())
	if err != nil {
		return fmt.Sprintf("error opening pod log stream, %v", err)
	}
	defer podStream.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podStream)
	if err != nil {
		return "error copying pod log stream to buf"
	}
	podLog = buf.String()

	return podLog
}

func getPodTabData(widgetLabelName string) (widgetNameLabel *widget.Label, widgetName *widget.Label, widgetNameScroll *container.Scroll) {
	widgetNameLabel = widget.NewLabel(widgetLabelName)
	widgetNameLabel.TextStyle = fyne.TextStyle{Monospace: true}
	widgetName = widget.NewLabel("")
	widgetName.TextStyle = fyne.TextStyle{Monospace: true}
	widgetName.Wrapping = fyne.TextWrapBreak
	widgetNameScroll = container.NewScroll(widgetName)
	widgetNameScroll.SetMinSize(fyne.Size{Height: 100})
	return widgetNameLabel, widgetName, widgetNameScroll
}

//TODO parse cluster context name to drop unnecessary text
func getCurrentContext() string {
	// get current context
	clientConfig, _ := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()
	return clientConfig.CurrentContext
}

// used by labels, annotations, ...
func convertMapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
