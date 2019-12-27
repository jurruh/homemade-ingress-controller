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

var hostTarget = map[string]string{}

type baseHandle struct{}

func (h *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	if target, ok := hostTarget[host]; ok {
		remoteUrl, err := url.Parse(target)
		if err != nil {
			log.Println("target parse fail:", err)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
		proxy.ServeHTTP(w, r)
		return
	}
	w.Write([]byte("403: Host forbidden " + host))
}

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
		backend := ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend
		hostTarget[ingress.Spec.Rules[0].Host] = "http://" + backend.ServiceName + "." + ingress.ObjectMeta.Namespace + ":" + backend.ServicePort.String()
	}

	fmt.Println(hostTarget)

	h := &baseHandle{}
	http.Handle("/", h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	log.Fatal(server.ListenAndServe())
}
