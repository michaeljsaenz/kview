package k8s

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/michaeljsaenz/kview/internal/utils"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

//TODO parse cluster context name to drop unnecessary text
func GetCurrentContext() string {
	// get current context
	clientConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()
	if err != nil {
		panic(err.Error())
	}
	return clientConfig.CurrentContext
}

func GetClientSet() *kubernetes.Clientset {
	// https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		//TODO raise this error to UI
		log.Fatal("kubeconfig error: ", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		//TODO raise this error to UI
		log.Fatal("clientset error: ", err)
	}

	return clientset

}

// get pod names to populate initial list
func GetPodData(c kubernetes.Clientset) (podData []string) {
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

func GetPodDetail(c kubernetes.Clientset, selectedPod string, podNamespace string) (string, string, string, string, string, string, []string) {
	pod, err := c.CoreV1().Pods(podNamespace).Get(context.TODO(), selectedPod, v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	podCreationTime := pod.GetCreationTimestamp()
	age := time.Since(podCreationTime.Time).Round(time.Second).String()
	ageHours := time.Since(podCreationTime.Time).Round(time.Hour).String()
	var ageHoursSlice []string
	if !strings.HasPrefix(ageHours, "0s") {
		ageHoursSlice = strings.Split(ageHours, "h")
		ageHoursInt, err := strconv.Atoi(ageHoursSlice[0])
		if err != nil {
			panic(err.Error())
		}
		if ageHoursInt > 23 {
			ageInt := ageHoursInt / 24
			age = strconv.Itoa(ageInt) + "d"
		}
	}

	var containers []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return string(pod.Status.Phase), age, string(pod.Namespace), utils.ConvertMapToString(pod.Labels),
		utils.ConvertMapToString(pod.Annotations), pod.Spec.NodeName, containers
}

func GetPodEvents(c kubernetes.Clientset, selectedPod string, podNamespace string) (podEvents []string) {
	events, _ := c.CoreV1().Events(podNamespace).List(context.TODO(), v1.ListOptions{FieldSelector: fmt.Sprintf("involvedObject.name=%s", selectedPod), TypeMeta: v1.TypeMeta{Kind: "Pod"}})
	for _, item := range events.Items {
		podEvents = append(podEvents, "~> "+item.EventTime.Time.Format("2006-01-02 15:04:05")+", "+item.Message)
	}
	return podEvents
}

func GetPodLogs(c kubernetes.Clientset, podNamespace string, selectedPod string, containerName string) (podLog string) {
	podLogReq := c.CoreV1().Pods(podNamespace).GetLogs(selectedPod, &corev1.PodLogOptions{Container: containerName,
		SinceSeconds: utils.CreateInt64(1800)}) //, LimitBytes: createInt64(1000)}) //TODO maybe grab less logs for display while usign refresh to pull new logs
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

func GetPodNamespace(c kubernetes.Clientset, podName string) (podNamespace string) {
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
