package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloudevents-explorer/internal/templates"
)

const gcsBaseURL = "http://localhost:4443/storage/v1"

type BucketsResponse struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Name string `json:"name"`
}

type ObjectsResponse struct {
	Items    []Object `json:"items,omitempty"`
	Prefixes []string `json:"prefixes,omitempty"`
}

type Object struct {
	Name string `json:"name"`
	Size int64  `json:"size,string"`
}

type ContentResponse struct {
	Content string `json:"content"`
}

func HandleGCS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, templates.GCS)
}

func HandleListBuckets(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(gcsBaseURL + "/b")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch buckets: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusInternalServerError)
		return
	}

	var gcsResp struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &gcsResp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse response: %v", err), http.StatusInternalServerError)
		return
	}

	buckets := make([]Bucket, len(gcsResp.Items))
	for i, item := range gcsResp.Items {
		buckets[i] = Bucket{Name: item.Name}
	}

	response := BucketsResponse{Buckets: buckets}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleListObjects(w http.ResponseWriter, r *http.Request) {
	bucket := r.URL.Query().Get("bucket")
	if bucket == "" {
		http.Error(w, "bucket parameter is required", http.StatusBadRequest)
		return
	}

	prefix := r.URL.Query().Get("prefix")

	url := fmt.Sprintf("%s/b/%s/o", gcsBaseURL, bucket)
	if prefix != "" {
		url += "?prefix=" + prefix + "&delimiter=/"
	} else {
		url += "?delimiter=/"
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch objects: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusInternalServerError)
		return
	}

	var gcsResp struct {
		Items []struct {
			Name string `json:"name"`
			Size string `json:"size"`
		} `json:"items"`
		Prefixes []string `json:"prefixes"`
	}

	if err := json.Unmarshal(body, &gcsResp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse response: %v", err), http.StatusInternalServerError)
		return
	}

	response := ObjectsResponse{
		Prefixes: gcsResp.Prefixes,
	}

	if len(gcsResp.Items) > 0 {
		objects := make([]Object, 0, len(gcsResp.Items))
		for _, item := range gcsResp.Items {
			var size int64
			fmt.Sscanf(item.Size, "%d", &size)
			objects = append(objects, Object{
				Name: item.Name,
				Size: size,
			})
		}
		response.Items = objects
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleGetObjectContent(w http.ResponseWriter, r *http.Request) {
	bucket := r.URL.Query().Get("bucket")
	object := r.URL.Query().Get("object")

	if bucket == "" || object == "" {
		http.Error(w, "bucket and object parameters are required", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/b/%s/o/%s?alt=media", gcsBaseURL, bucket, strings.ReplaceAll(object, "/", "%2F"))

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch object: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusInternalServerError)
		return
	}

	response := ContentResponse{Content: string(body)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleDownloadObject(w http.ResponseWriter, r *http.Request) {
	bucket := r.URL.Query().Get("bucket")
	object := r.URL.Query().Get("object")

	if bucket == "" || object == "" {
		http.Error(w, "bucket and object parameters are required", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/b/%s/o/%s?alt=media", gcsBaseURL, bucket, strings.ReplaceAll(object, "/", "%2F"))

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch object: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set headers for file download
	filename := object
	if idx := strings.LastIndex(object, "/"); idx != -1 {
		filename = object[idx+1:]
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copy the content to response
	io.Copy(w, resp.Body)
}
