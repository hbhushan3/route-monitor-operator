package syntheticmonitoring

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	// "github.com/openshift/hypershift/api/v1beta1"
	// "github.com/openshift/route-monitor-operator/api/v1alpha1"
	// "github.com/openshift/route-monitor-operator/pkg/util/finalizer"
	// utilreconcile "github.com/openshift/route-monitor-operator/pkg/util/reconcile"

	// v1 "k8s.io/api/core/v1"
	// kerr "k8s.io/apimachinery/pkg/api/errors" can use for error catching
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" can use for martshalling, datat types
	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	// "sigs.k8s.io/controller-runtime/pkg/builder"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/handler"
	// "sigs.k8s.io/controller-runtime/pkg/manager"
	// "sigs.k8s.io/controller-runtime/pkg/predicate"
	// "sigs.k8s.io/controller-runtime/pkg/source"

	"bytes"
	"encoding/json"
	"net/http"
)

var logger logr.Logger = ctrl.Log.WithName("controllers").WithName("HostedControlPlane/SyntheticMonitoring")

type APIClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

func NewAPIClient(baseURL, apiToken string) *APIClient {
	return &APIClient{
		baseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: &http.Client{},
	}
}

func (client *APIClient) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	url := client.baseURL + path
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Api-Token "+client.apiToken)
	req.Header.Set("Content-Type", "application/json")

	return client.httpClient.Do(req)
}

func (client *APIClient) createDynatraceHTTPMonitor(ctx context.Context, monitorName, url string) (string, error) {
	monitorConfig := map[string]interface{}{
		// "name":         monitorName,
		// "type":         "http",
		// "frequencyMin": 5,
		// "enabled":      true,
		// "locations":    []string{"GEOLOCATION-9999453BE4BDB3CD"},
		// "config": map[string]string{
		// 	"url": url,
		// },
		"name":         monitorName,
		"frequencyMin": 5,
		"enabled":      true,
		"type":         "HTTP",
		"script": map[string]interface{}{
			"version": "1.0",
			"requests": []map[string]interface{}{
				{
					"description": "api availability",
					"url":         url,
					"method":      "GET",
					"requestBody": "",
					"configuration": map[string]interface{}{
						"acceptAnyCertificate": true,
						"followRedirects":      true,
					},
					"preProcessingScript":  "",
					"postProcessingScript": "",
				},
			},
		},
		"locations": []string{"GEOLOCATION-9999453BE4BDB3CD"},
		"anomalyDetection": map[string]interface{}{
			"outageHandling": map[string]interface{}{
				"globalOutage": true,
				"localOutage":  false,
				"localOutagePolicy": map[string]interface{}{
					"affectedLocations": 1,
					"consecutiveRuns":   3,
				},
			},
			"loadingTimeThresholds": map[string]interface{}{
				"enabled": false,
				"thresholds": []map[string]interface{}{
					{
						"type":    "TOTAL",
						"valueMs": 60000,
					},
				},
			},
		},
		"tags": []map[string]interface{}{
			{
				"key":   "syntheticMonitorEnabled",
				"value": "true",
			},
		},
	}

	resp, err := client.makeRequest("POST", "/synthetic/monitors", monitorConfig)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to create HTTP monitor. Status code: %d", resp.StatusCode)
	}
	fmt.Println("HTTP monitor created successfully.")

	var createdMonitor map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&createdMonitor)
	if err != nil {
		return "", err
	}
	monitorID := createdMonitor["entityId"].(string)
	fmt.Printf("Monitor ID is %s \n", monitorID)
	return monitorID, nil
}

func (client *APIClient) deleteDynatraceHTTPMonitor(ctx context.Context, monitorID string) error {
	path := fmt.Sprintf("/synthetic/monitors/%s", monitorID)

	resp, err := client.makeRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to delete monitor. Status code: %d", resp.StatusCode)
	}

	fmt.Printf("Successfully deleted monitor with ID %s \n", monitorID)
	return nil
}
