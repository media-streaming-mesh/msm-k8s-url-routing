package url_handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func createMockDeployment(dep *v1.Deployment, clientset kubernetes.Interface, ns string) {
	client := clientset.AppsV1().Deployments(ns)

	_, err := client.Create(context.TODO(), dep, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("error creating mock deployment")
	}
}

func createMockService(svc *apiv1.Service, clientset kubernetes.Interface, ns string) {
	client := clientset.CoreV1().Services(ns)

	_, err := client.Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("error creating mock service")
	}
}

func TestGetInternalURLs(t *testing.T) {
	urlTests := map[string][]string{
		"http://localhost:9898/apiv1/url-routing?url=http://kubernetes:443":                          []string{"http://"},
		"http://localhost:9898/apiv1/url-routing?url=http://http://localhost:9898/apiv1/url-routing": []string{"9898/apiv1/url-routing"},
	}

	for url, expectedResp := range urlTests {
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("error connecting to svc: %v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error reading response: %v", err)
		}

		var response []string
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatalf("error unmarshaling response: %v", err)
		}

		if !strings.Contains(response[0], expectedResp[0]) {
			t.Fatalf("responses do not match: %s %s", response[0], expectedResp[0])
		}
	}
}
