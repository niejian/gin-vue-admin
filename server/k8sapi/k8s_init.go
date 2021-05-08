// k8sapi doc

package k8sapi

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

var clientSet *kubernetes.Clientset

func init() {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)

	clientSetInit, err := kubernetes.NewForConfig(config)

	if nil != clientSetInit {
		clientSet = clientSetInit
	} else {
		panic("k8s connect failed")
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Println("connect k8s success")
	}
}

func InitK8s() *kubernetes.Clientset {

	return clientSet
}
