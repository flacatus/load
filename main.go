package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	appstudioApi "github.com/redhat-appstudio/application-api/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	scheme = runtime.NewScheme()
)

type Concurently struct {
	KubernetesClient crclient.Client
}

// Return a rest client to perform CRUD operations on Kubernetes objects
func (c *Concurently) KubeRest() crclient.Client {
	return c.KubernetesClient
}

func NewConcurentlyController() (*Concurently, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	crClient, err := crclient.New(cfg, crclient.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	return &Concurently{
		KubernetesClient: crClient,
	}, nil
}

func (c *Concurently) CreateApplication(name string, namespace string, wg *sync.WaitGroup) (*appstudioApi.Application, error) {
	defer wg.Done()
	fmt.Println(name)
	application := &appstudioApi.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appstudioApi.ApplicationSpec{
			DisplayName: name,
		},
	}
	err := c.KubeRest().Create(context.TODO(), application)
	if err != nil {
		return nil, err
	}

	if err := WaitUntil(c.ApplicationDevfilePresent(application), time.Second*30); err != nil {
		return nil, fmt.Errorf("timed out when waiting for devfile content creation for application %s in %s namespace: %+v", name, namespace, err)
	}

	return application, nil
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appstudioApi.AddToScheme(scheme))
}

func main() {
	controller, _ := NewConcurentlyController()
	var wg sync.WaitGroup

	wg.Add(2)
	go controller.CreateApplication("app1", "user1", &wg)
	go controller.CreateApplication("app2", "user1", &wg)
	go controller.CreateApplication("app3", "user1", &wg)
	go controller.CreateApplication("app4", "user1", &wg)
	go controller.CreateApplication("app5", "user1", &wg)
	go controller.CreateApplication("app6", "user1", &wg)
	go controller.CreateApplication("app7", "user1", &wg)
	go controller.CreateApplication("app8", "user1", &wg)
	go controller.CreateApplication("app9", "user1", &wg)
	go controller.CreateApplication("app10", "user1", &wg)
	go controller.CreateApplication("app11", "user1", &wg)
	go controller.CreateApplication("app12", "user1", &wg)

	wg.Wait()
	fmt.Println("Done!")
}
