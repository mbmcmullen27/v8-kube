package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"path/filepath"
	// "time"
	"io/ioutil"
	"encoding/json"
	
	
	"rogchap.com/v8go"

	// "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func configure() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func execute(pod string, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	dat, _ := ioutil.ReadFile("util.js")
	ctx, _ := v8go.NewContext() 
	
	ctx.RunScript(string(dat), "util.js") 
	var scr string ="const result = parse("+pod+")"
	ctx.RunScript(scr, "main.js") 
	val, _ := ctx.RunScript("result", "value.js") 
	
	fmt.Printf("%s : %s\n",name, val)
}

func main() {
	clientset := configure()
	pods, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})

	fmt.Printf("%d\n",len(pods.Items))
	fmt.Printf("%T\n", pods)
	
	length := len(pods.Items)

	var wg sync.WaitGroup
	
	for i:=0; i<length; i++ {
		if err != nil {
			panic(err.Error())
		}

		name:=pods.Items[i].ObjectMeta.Name
		data, _ := json.Marshal(pods.Items[i])

		wg.Add(1)
		go execute(string(data), name, &wg)
	}

	wg.Wait()
}