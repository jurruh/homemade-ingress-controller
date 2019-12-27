package main

import (
	"fmt"
	"net/http"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	fmt.Println("Starting Ingress Controller!")

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("No in cluster config trying to use ~/.kube/config")
		home, _ := os.UserHomeDir()
		config, err = clientcmd.BuildConfigFromFlags("", home+"/.kube/config")

		if err != nil {
			panic(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ingresses, _ := clientset.NetworkingV1beta1().Ingresses("").List(v1.ListOptions{})

	fmt.Print(ingresses)

	http.ListenAndServe(":8080", nil)
}
