package k8s

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/michaeljsaenz/kview/internal/utils"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
)

// TODO parse cluster context name to drop unnecessary text
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

// get pod names with provided namespace
func GetPodDataWithNamespace(c kubernetes.Clientset, namespace string) (podData []string) {
	pods, err := c.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		podData = append(podData, fmt.Sprint(err))
	} else {
		for _, pod := range pods.Items {
			podData = append(podData, pod.Name)
		}
	}
	return podData

}

// get namespaces
func GetNamespaces(c kubernetes.Clientset) (namespaceList []string) {
	// retrieve the list of namespaces
	namespaces, err := c.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Printf("failed to get namespaces: %v", err)
		namespaceList = append(namespaceList, fmt.Sprint(err))
	} else {
		for _, namespace := range namespaces.Items {
			namespaceList = append(namespaceList, namespace.Name)
		}
	}
	return namespaceList

}

func GetPodDetail(c kubernetes.Clientset, selectedPod string, podNamespace string) (string, string, string, []string) {
	pod, err := c.CoreV1().Pods(podNamespace).Get(context.TODO(), selectedPod, v1.GetOptions{})
	if err != nil {
		fmt.Printf("failed to get pod detail: %v", err)
	}

	podCreationTime := pod.GetCreationTimestamp()
	age := time.Since(podCreationTime.Time).Round(time.Second)
	podAge := age.String()
	if int(math.Trunc(age.Hours())) >= 24 {
		ageInDays := int(math.Trunc(age.Hours())) / 24
		podAge = strconv.Itoa(ageInDays) + "d"
	}

	var containers []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return string(pod.Status.Phase), podAge, pod.Spec.NodeName, containers
}

func GetPodLabels(c kubernetes.Clientset, selectedPod string, podNamespace string) string {
	pod, err := c.CoreV1().Pods(podNamespace).Get(context.TODO(), selectedPod, v1.GetOptions{})
	if err != nil {
		fmt.Printf("failed to get pod labels: %v", err)
	}

	return utils.ConvertMapToString(pod.Labels)
}

func GetPodAnnotations(c kubernetes.Clientset, selectedPod string, podNamespace string) string {
	pod, err := c.CoreV1().Pods(podNamespace).Get(context.TODO(), selectedPod, v1.GetOptions{})
	if err != nil {
		fmt.Printf("failed to get pod annotations: %v", err)
	}

	return utils.ConvertMapToString(pod.Annotations)
}

func GetPodEvents(c kubernetes.Clientset, selectedPod string, podNamespace string) (podEvents []string) {
	events, _ := c.CoreV1().Events(podNamespace).List(context.TODO(), v1.ListOptions{FieldSelector: fmt.Sprintf("involvedObject.name=%s", selectedPod), TypeMeta: v1.TypeMeta{Kind: "Pod"}})
	for _, item := range events.Items {
		podEvents = append(podEvents, item.FirstTimestamp.String()+" "+item.Message)
	}
	return podEvents
}

func GetPodLogs(c kubernetes.Clientset, podNamespace string, selectedPod string, containerName string) (podLog string) {
	const (
		logTailLines = 1000
	)
	podLogReq := c.CoreV1().Pods(podNamespace).GetLogs(selectedPod, &corev1.PodLogOptions{Container: containerName,
		TailLines: utils.CreateInt64(logTailLines)})

	podLog = podLogStreamToString(podLogReq)

	return podLog
}

func podLogStreamToString(podLogReq *rest.Request) (podLog string) {
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

func GetPodYaml(c kubernetes.Clientset, podNamespace string, podName string) (string, error) {
	pod, err := c.CoreV1().Pods(podNamespace).Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error getting pod: %v", err)
	}

	// clear unnecessary fields
	pod.ObjectMeta.ManagedFields = nil
	pod.ObjectMeta.GenerateName = ""
	pod.Status = corev1.PodStatus{}

	// serialize the Pod to YAML format
	codec := serializer.NewCodecFactory(scheme.Scheme).LegacyCodec(corev1.SchemeGroupVersion)
	marshaledYaml, err := runtime.Encode(codec, pod)
	if err != nil {
		return "", fmt.Errorf("error encoding YAML: %v", err)
	}

	// convert the marshaled YAML to a string
	yamlString, err := yaml.JSONToYAML(marshaledYaml)
	if err != nil {
		return "", fmt.Errorf("error converting YAML to string: %v", err)
	}

	return string(yamlString), nil

}
