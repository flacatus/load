package main

import (
	"context"
	"time"

	appstudioApi "github.com/redhat-appstudio/application-api/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (c *Concurently) ApplicationDevfilePresent(application *appstudioApi.Application) wait.ConditionFunc {
	return func() (bool, error) {
		app, err := c.GetHasApplication(application.Name, application.Namespace)
		if err != nil {
			return false, nil
		}
		application.Status = app.Status
		return application.Status.Devfile != "", nil
	}
}

// GetHasApplication return the Application Custom Resource object
func (c *Concurently) GetHasApplication(name, namespace string) (*appstudioApi.Application, error) {
	namespacedName := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	application := appstudioApi.Application{
		Spec: appstudioApi.ApplicationSpec{},
	}
	err := c.KubeRest().Get(context.TODO(), namespacedName, &application)
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func WaitUntil(cond wait.ConditionFunc, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, cond)
}
