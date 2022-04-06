package pkg

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

func Stash(client kubernetes.Interface, namespace string) error {
	deploys, err := client.AppsV1().Deployments(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, d := range deploys.Items {
		replicas := d.Spec.Replicas

		if *replicas == int32(0) {
			fmt.Printf("\"%s\" replicas is 0, don't need stash \n", d.Name)
			continue
		}

		d.Spec.Replicas = int32Ptr(0)
		d.Annotations["stash-replicas"] = strconv.Itoa(int(*replicas))

		_, err = client.AppsV1().Deployments(namespace).Update(context.Background(), &d, v1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("\"%s\" stash replicas succeed \n", d.Name)
	}
	return nil
}

func Recover(client kubernetes.Interface, namespace string) error {
	deploys, err := client.AppsV1().Deployments(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, d := range deploys.Items {
		if replicas, ok := extractReplicas(d.ObjectMeta); ok {
			d.Spec.Replicas = replicas
			_, err = client.AppsV1().Deployments(namespace).Update(context.Background(), &d, v1.UpdateOptions{})
			if err != nil {
				return err
			}
			fmt.Printf("\"%s\" recover replicas %d succeed \n", d.Name, *d.Spec.Replicas)
		} else {
			fmt.Printf("\"%s\" no stash,  don't need recover! \n", d.Name)
		}
	}
	return nil
}

func extractReplicas(meta v1.ObjectMeta) (*int32, bool) {
	stashReplicas := meta.Annotations["stash-replicas"]
	meta.Annotations["stash-replicas"] = ""
	if stashReplicas == "" {
		return int32Ptr(0), false
	}
	replicas, err := strconv.Atoi(stashReplicas)
	if err != nil {
		return int32Ptr(0), false
	}
	return int32Ptr(int32(replicas)), true
}

func int32Ptr(i int32) *int32 { return &i }
