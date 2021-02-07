package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"

	"github.com/mredvard/kubernetes-crd-example/api/types/v1alpha1"
	clientV1alpha1 "github.com/mredvard/kubernetes-crd-example/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("config", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("config-path", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		panic(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	projects, err := clientSet.Projects("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("projects found: %+v\n", projects)

	store := WatchResources(clientSet)

	for {
		projectsFromStore := store.List()
		fmt.Printf("project in store: %d\n", len(projectsFromStore))

		time.Sleep(2 * time.Second)
	}
}
