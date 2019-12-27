package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("Starting Ingress Controller!")

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

	for _, ingress := range ingresses.Items {
		fmt.Println(ingress.Spec.Rules[0].Host)
	}

	fmt.Println("Start http")

	targetURL := "https://google.com"

	u, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Could not parse downstream url: %s", targetURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = req.URL.Host
	}

	http.HandleFunc("/", proxy.ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
