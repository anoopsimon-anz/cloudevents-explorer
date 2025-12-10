package handlers

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloudevents-explorer/internal/templates"
)

type RestRequest struct {
	Method  string                 `json:"method"`
	URL     string                 `json:"url"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"`
	TLSCert string                 `json:"tlsCert"`
	TLSKey  string                 `json:"tlsKey"`
}

type RestResponse struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]string      `json:"headers"`
	Body       interface{}            `json:"body"`
	Error      string                 `json:"error,omitempty"`
}

func HandleRestClient(w http.ResponseWriter, r *http.Request) {
	html := templates.GetBaseHTML("REST Client", templates.RestClient, templates.RestClientJS)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func HandleRestSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendRestError(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create HTTP client with TLS configuration
	// Default to InsecureSkipVerify=true for testing convenience
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Configure TLS if certificates are provided
	if req.TLSCert != "" && req.TLSKey != "" {
		cert, err := tls.X509KeyPair([]byte(req.TLSCert), []byte(req.TLSKey))
		if err != nil {
			sendRestError(w, "Invalid TLS certificate or key: "+err.Error(), http.StatusBadRequest)
			return
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.RootCAs = x509.NewCertPool()
	} else if req.TLSCert != "" {
		// If only cert is provided, use it as CA cert
		tlsConfig.RootCAs = x509.NewCertPool()
		if ok := tlsConfig.RootCAs.AppendCertsFromPEM([]byte(req.TLSCert)); !ok {
			sendRestError(w, "Failed to parse TLS certificate", http.StatusBadRequest)
			return
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Prepare request body
	var bodyReader io.Reader
	if req.Body != nil && len(req.Body) > 0 {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			sendRestError(w, "Failed to marshal request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		sendRestError(w, "Failed to create request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	resp, err := client.Do(httpReq)
	if err != nil {
		sendRestError(w, "Request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		sendRestError(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse response body as JSON if possible
	var parsedBody interface{}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &parsedBody); err != nil {
			// If not JSON, return as string
			parsedBody = string(respBody)
		}
	}

	// Extract response headers
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// Send successful response
	response := RestResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       parsedBody,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendRestError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(RestResponse{
		StatusCode: statusCode,
		Error:      message,
	})
}
