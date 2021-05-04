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

func execute(pod string, wg *sync.WaitGroup, file *os.File) {
	defer wg.Done()

	//currently runs 1 isolate per pod
	iso, _ := v8go.NewIsolate()

	//callback to write a byteslice to the opened flie
	filewrite, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Printf("%v", info.Args()[0])
		data := []byte(info.Args()[0].String())
		numb, err := file.Write(data)
		check(err)
		fmt.Printf("wrote %d bytes\n", numb)

		file.Sync()

		return nil
	})

	global, _ := v8go.NewObjectTemplate(iso)
	global.Set("print", filewrite)

	util, _ := ioutil.ReadFile("util.js")
	ctx, _ := v8go.NewContext(iso, global) 
	
	ctx.RunScript(string(util), "util.js") 
	var scr string ="const result = parse("+pod+")"
	ctx.RunScript(scr, "main.js") 
	ctx.RunScript("result", "value.js") 

}

func main() {
	clientset := configure()
	pods, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})

	length := len(pods.Items)
	fmt.Printf("Found %d pods\n",length)

	//open an output data file 
	file, err := os.Create("/tmp/kubdat") 
    check(err)
	defer file.Close()

	var wg sync.WaitGroup
	
	for i:=0; i<length; i++ {
		if err != nil {
			panic(err.Error())
		}

		name:=pods.Items[i].ObjectMeta.Name
		data, _ := json.Marshal(pods.Items[i])

		wg.Add(1)
		go execute(string(data), name, &wg, file)
	}

	wg.Wait()
}