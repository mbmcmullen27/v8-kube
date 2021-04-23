package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"path/filepath"
	// "time"
	"encoding/json"
	
	
	"rogchap.com/v8go"

	// "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
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
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	
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
		// result:= execute(string(data))
		wg.Add(1)
		go execute(string(data), name, &wg)
		// fmt.Printf("%s : %s\n",name, result)
	}

	wg.Wait()
}

func execute(pod string, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, _ := v8go.NewContext() // creates a new V8 context with a new Isolate aka VM
	ctx.RunScript("const log = (a) => Object.keys(a)", "log.js") // executes a script on the global context
	var scr string ="const result = log("+pod+")"
	ctx.RunScript(scr, "main.js") // any functions previously added to the context can be called
	val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
	fmt.Printf("%s : %s\n",name, val)
}

