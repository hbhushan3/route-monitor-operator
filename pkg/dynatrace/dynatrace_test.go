package hostedcontrolplane

import (
	"net/http"
	"net/http/httptest"
	"testing"

	hypershiftv1beta1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
)

func setupMockServer(handlerFunc http.HandlerFunc) string {
	mockServer := httptest.NewServer(handlerFunc)
	return mockServer.URL
}
func createMockHandlerFunc(responseBody string, statusCode int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	})
}
func TestNewAPIClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		apiToken string
	}{
		{
			name:     "Valid API Client Initialization",
			baseURL:  "https://example.com/api",
			apiToken: "mockToken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiClient := NewDynatraceAPIClient(tt.baseURL, tt.apiToken)

			if apiClient.baseURL != tt.baseURL {
				t.Errorf("Expected baseURL to be %s, got %s", tt.baseURL, apiClient.baseURL)
			}

			if apiClient.apiToken != tt.apiToken {
				t.Errorf("Expected apiToken to be %s, got %s", tt.apiToken, apiClient.apiToken)
			}

			if apiClient.httpClient == nil {
				t.Error("Expected httpClient to be initialized, got nil")
			}
		})
	}
}

func TestAPIClient_makeRequest(t *testing.T) {
	// Define test cases in a slice of structs
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "Make a GET request",
			method:         "GET",
			body:           "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Make a POST request",
			method:         "POST",
			body:           `{"key": "value"}`,
			expectedStatus: http.StatusOK,
		},
	}

	// Iterate through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a custom handler function for the mock server
			handlerFunc := createMockHandlerFunc(`{"message": "Mocked response"}`, http.StatusOK)
			// Create the mock server
			mockServerURL := setupMockServer(handlerFunc)
			// Create an instance of the APIClient
			apiClient := NewDynatraceAPIClient(mockServerURL, "mockedToken")

			// Make the request
			response, err := apiClient.MakeRequest(tt.method, "/test", tt.body)
			if err != nil {
				t.Errorf("Error making %s request: %v", tt.method, err)
			}
			defer response.Body.Close()

			// Assert the response status code
			if response.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, response.StatusCode)
			}
		})
	}
}

func TestAPIClient_GetDynatraceEquivalentClusterRegionId(t *testing.T) {
	tests := []struct {
		name               string
		region             string
		hostedControlPlane *hypershiftv1beta1.HostedControlPlane
		mockResponse       string
		mockStatusCode     int
		expectedRegionID   string
		expectedError      string
	}{
		{
			name:             "Valid region code and location found",
			region:           "us-west-2",
			mockResponse:     `{"locations": [{"name": "Oregon", "type": "PUBLIC", "cloudPlatform": "AMAZON_EC2", "entityId": "123"}]}`,
			mockStatusCode:   http.StatusOK,
			expectedRegionID: "123",
			expectedError:    "",
		},
		{
			name:             "Invalid region code (no matching location)",
			region:           "invalid-region-code",
			mockResponse:     `{"locations": []}`,
			mockStatusCode:   http.StatusBadRequest,
			expectedRegionID: "",
			expectedError:    "location not found for region: invalid-region-code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mocking the HTTP server to return the desired response
			handlerFunc := createMockHandlerFunc(tt.mockResponse, tt.mockStatusCode)
			// Create the mock server using the setupMockServer function
			mockServer := setupMockServer(handlerFunc)
			// Create an instance of the APIClient using the reusable setup
			mockClient := NewDynatraceAPIClient(mockServer, "mockedToken")

			// Call the function being tested
			regionID, err := mockClient.GetDynatraceEquivalentClusterRegionId(tt.region)

			// Check the returned values against the expected results
			if regionID != tt.expectedRegionID {
				t.Errorf("Got: %s, Expected: %s", regionID, tt.expectedRegionID)
			}

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("Got error: %v, Expected error: %s", err, tt.expectedError)
			}
		})
	}
}

func TestAPIClient_CreateDynatraceHTTPMonitor(t *testing.T) {
	// Mocked response data for testing
	mockMonitorName := "TestMonitor"
	mockApiUrl := "https://example.com"
	mockClusterId := "12345"
	mockDynatraceEquivalentClusterRegionId := "us-east-1"

	// Create a list of test cases
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		expectedID     string
		expectError    bool
	}{
		{
			name:           "SuccessfulMonitorCreation",
			mockResponse:   `{"entityId": "56789"}`,
			mockStatusCode: http.StatusOK,
			expectedID:     "56789",
			expectError:    false,
		},
		{
			name:           "ErrorResponseFromServer",
			mockResponse:   "Bad request",
			mockStatusCode: http.StatusBadRequest,
			expectedID:     "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the HTTP server to return the desired response
			mockServer := setupMockServer(createMockHandlerFunc(tt.mockResponse, tt.mockStatusCode))

			// Create an instance of the APIClient using the mock server
			mockClient := NewDynatraceAPIClient(mockServer, "mockedToken")

			// Call the method under test
			monitorID, err := mockClient.CreateDynatraceHTTPMonitor(mockMonitorName, mockApiUrl, mockClusterId, mockDynatraceEquivalentClusterRegionId)

			// Check for errors or expected values based on the test case
			if (err != nil) != tt.expectError {
				t.Errorf("Unexpected error status. Expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError && monitorID != tt.expectedID {
				t.Errorf("Incorrect monitor ID. Expected: %s, Got: %s", tt.expectedID, monitorID)
			}
		})
	}
}

func TestAPIClient_DeleteDynatraceHTTPMonitor(t *testing.T) {
	// Create a list of test cases
	tests := []struct {
		name           string
		mockStatusCode int
		expectError    bool
	}{
		{
			name:           "Successful DELETE request",
			mockStatusCode: http.StatusNoContent,
			expectError:    false,
		},
		{
			name:           "Failed DELETE request",
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the HTTP server to return the desired response
			mockServer := setupMockServer(createMockHandlerFunc("", tt.mockStatusCode))
			apiClient := NewDynatraceAPIClient(mockServer, "mockedToken")

			// Call the method under test
			err := apiClient.DeleteDynatraceHTTPMonitor("123")

			// Check for errors based on the expected outcome
			if (err != nil) != tt.expectError {
				t.Errorf("Unexpected error status. Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}