package main

import (
	"log"
	"net/http"

	"cloudevents-explorer/internal/config"
	"cloudevents-explorer/internal/handlers"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
	}

	// Register page handlers
	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/pubsub", handlers.HandlePubSub)
	http.HandleFunc("/kafka", handlers.HandleKafka)
	http.HandleFunc("/flow-diagram", handlers.HandleFlowDiagram)

	// Register API handlers
	http.HandleFunc("/api/configs", handlers.HandleGetConfigs)
	http.HandleFunc("/api/pubsub/configs", handlers.HandleSavePubSubConfig)
	http.HandleFunc("/api/kafka/configs", handlers.HandleSaveKafkaConfig)
	http.HandleFunc("/api/pubsub/pull", handlers.HandlePullPubSub)
	http.HandleFunc("/api/kafka/pull", handlers.HandlePullKafka)
	http.HandleFunc("/api/kafka/publish", handlers.HandlePublishKafka)

	port := "8888"
	log.Printf("🚀 CloudEvents Explorer starting on http://localhost:%s", port)
	log.Printf("📝 Configuration file: configs.json")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
