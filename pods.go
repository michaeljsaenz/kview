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
)

// get pod names to populate initial list
func getPodData(c kubernetes.Clientset) (podData []string) {
	pods, err := c.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		podData = append(podData, fmt.Sprint(err))
	} else {
		for _, pod := range pods.Items {
			podData = append(podData, pod.Name)
		}
	}
	return podData
}

func getPodDetail(c kubernetes.Clientset, selectedPod string, podNamespace string) (string, string, string, string, string, string, []string) {
	pod, err := c.CoreV1().Pods(podNamespace).Get(context.TODO(), selectedPod, v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	podCreationTime := pod.GetCreationTimestamp()
	age := time.Since(podCreationTime.Time).Round(time.Second)
	var containers []string

	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return string(pod.Status.Phase), age.String(), string(pod.Namespace), convertMapToString(pod.Labels),
		convertMapToString(pod.Annotations), pod.Spec.NodeName, containers
}

func getPodEvents(c kubernetes.Clientset, selectedPod string, podNamespace string) (podEvents []string) {
	events, _ := c.CoreV1().Events(podNamespace).List(context.TODO(), v1.ListOptions{FieldSelector: fmt.Sprintf("involvedObject.name=%s", selectedPod), TypeMeta: v1.TypeMeta{Kind: "Pod"}})
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

func getPodNamespace(c kubernetes.Clientset, podName string) (podNamespace string) {
	podNameWithNamespace := make(map[string]string)
	pods, err := c.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var podsItemsList []string
	for _, pod := range pods.Items {
		podsItemsList = append(podsItemsList, pod.Name, pod.Namespace)
	}
	for i := 0; i < len(podsItemsList); i += 2 {
		podNameWithNamespace[podsItemsList[i]] = podsItemsList[i+1]
	}
	podNamespace = podNameWithNamespace[podName]

	return podNamespace
}
