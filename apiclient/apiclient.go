package apiclient

type ApiClient interface{}

type apiClientImpl struct {
}

func New() ApiClient {
	return &apiClientImpl{}
}

// curl "http://localhost:8080/h1?p1=abc"
// curl -X POST -H "Content-Type: application/json" -d '{"r1":"xyz", "r2":"abc"}' "http://localhost:8080/h2"

//func callExternalAPI() (*ExternalAPIResponse, error) {
//	// Create HTTP client with timeout
//	client := &http.Client{
//		Timeout: 5 * time.Second,
//	}
//
//	// Create request
//	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
//	if err != nil {
//		return nil, fmt.Errorf("error creating request: %w", err)
//	}
//
//	// Set headers if needed
//	req.Header.Set("Accept", "application/json")
//
//	// Send request
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("error sending request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// Check status code
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//
//	var externalResp RequestObject
//	if err := json.NewDecoder(resp.Body).Decode(&externalResp); err != nil {
//		errorLogger.Println("failed to decode request body:", err)
//		http.Error(w, "Bad Request", http.StatusBadRequest)
//		return
//	}
//
//	return &externalResp, nil
//}
