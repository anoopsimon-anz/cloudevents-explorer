package kafka

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/linkedin/goavro/v2"

	"cloudevents-explorer/internal/types"
)

type PullParams struct {
	Brokers        string `json:"brokers"`
	Topic          string `json:"topic"`
	ConsumerGroup  string `json:"consumerGroup"`
	SchemaRegistry string `json:"schemaRegistry"`
	MaxMessages    int    `json:"maxMessages"`
}

type PullResult struct {
	Messages []types.CloudEvent `json:"messages"`
	Count    int                `json:"count"`
}

type PublishParams struct {
	Brokers        string                 `json:"brokers"`
	Topic          string                 `json:"topic"`
	SchemaRegistry string                 `json:"schemaRegistry"`
	Message        map[string]interface{} `json:"message"`
}

type PublishResult struct {
	Status    string `json:"status"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
}

func decodeAvroMessage(data []byte, schemaRegistryURL string) (map[string]interface{}, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("message too short")
	}

	// First byte is magic byte (should be 0)
	// Next 4 bytes are schema ID (big-endian)
	schemaID := binary.BigEndian.Uint32(data[1:5])

	// Fetch schema from registry
	schemaURL := fmt.Sprintf("%s/schemas/ids/%d", schemaRegistryURL, schemaID)
	resp, err := http.Get(schemaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("schema registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema response: %w", err)
	}

	var schemaResp struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(body, &schemaResp); err != nil {
		return nil, fmt.Errorf("failed to parse schema response: %w", err)
	}

	// Create Avro codec
	codec, err := goavro.NewCodec(schemaResp.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create codec: %w", err)
	}

	// Decode the message (skip first 5 bytes)
	native, _, err := codec.NativeFromBinary(data[5:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode avro: %w", err)
	}

	// Convert to map[string]interface{}
	if result, ok := native.(map[string]interface{}); ok {
		return result, nil
	}

	return nil, fmt.Errorf("decoded data is not a map")
}

func Pull(params PullParams) (*PullResult, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": params.Brokers,
		"group.id":          params.ConsumerGroup,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	defer c.Close()

	err = c.Subscribe(params.Topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	messages := []types.CloudEvent{}

	// Use a more aggressive timeout strategy:
	// - Overall timeout: 10 seconds (enough time to fetch messages)
	// - No-message timeout: 2 seconds of consecutive failures before giving up
	timeout := time.After(10 * time.Second)
	noMessageTimeout := 2 * time.Second
	lastMessageTime := time.Now()

	for len(messages) < params.MaxMessages {
		select {
		case <-timeout:
			goto done
		default:
			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// If we haven't seen a message in noMessageTimeout, and we have at least 1 message, stop
				if len(messages) > 0 && time.Since(lastMessageTime) > noMessageTimeout {
					goto done
				}
				continue
			}

			// We got a message! Update the last message time
			lastMessageTime = time.Now()

			event := types.CloudEvent{
				ID:        fmt.Sprintf("%s-%d-%d", params.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset),
				Subject:   params.Topic,
				Published: msg.Timestamp.Format(time.RFC3339),
				Timestamp: msg.Timestamp.Unix(),
			}

			// Decode Avro message using schema registry
			if params.SchemaRegistry != "" && len(msg.Value) > 5 {
				decodedData, err := decodeAvroMessage(msg.Value, params.SchemaRegistry)
				if err == nil && decodedData != nil {
					event.Data = decodedData
				} else {
					// Log the error for debugging
					fmt.Printf("Failed to decode Avro message (partition=%d, offset=%d): %v\n",
						msg.TopicPartition.Partition, msg.TopicPartition.Offset, err)

					// Extract schema ID from message
					var schemaID uint32
					if len(msg.Value) >= 5 {
						schemaID = binary.BigEndian.Uint32(msg.Value[1:5])
					}

					// Store error information with helpful debugging details
					event.Data = map[string]interface{}{
						"_error":          "Avro Decoding Failed",
						"_errorDetails":   err.Error(),
						"_schemaRegistry": params.SchemaRegistry,
						"_schemaID":       schemaID,
						"_messageSize":    len(msg.Value),
						"_troubleshooting": map[string]string{
							"step1": fmt.Sprintf("Verify schema registry is accessible: %s", params.SchemaRegistry),
							"step2": fmt.Sprintf("Check if schema ID %d exists: %s/schemas/ids/%d", schemaID, params.SchemaRegistry, schemaID),
							"step3": "Verify container can reach dep_redpanda:18081",
							"step4": "Check if running with correct config (Docker vs Local)",
						},
					}
				}
			} else {
				// Try plain JSON first
				var data map[string]interface{}
				if err := json.Unmarshal(msg.Value, &data); err == nil {
					event.Data = data
				} else {
					event.RawData = string(msg.Value)
				}
			}

			messages = append(messages, event)
		}
	}

done:
	// Reverse messages so newest appears first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return &PullResult{
		Messages: messages,
		Count:    len(messages),
	}, nil
}

func Publish(params PublishParams) (*PublishResult, error) {
	// Convert message to JSON bytes
	messageJSON, err := json.Marshal(params.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": params.Brokers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	defer p.Close()

	// Encode to Avro if schema registry is configured
	var messageBytes []byte
	if params.SchemaRegistry != "" {
		// Fetch schema from registry (using subject for Unica events)
		schemaURL := fmt.Sprintf("%s/subjects/au.data.unica.comms.event/versions/latest", params.SchemaRegistry)
		resp, err := http.Get(schemaURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch schema: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("schema registry returned status %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read schema response: %w", err)
		}

		var schemaResp struct {
			Schema string `json:"schema"`
			ID     int    `json:"id"`
		}
		if err := json.Unmarshal(body, &schemaResp); err != nil {
			return nil, fmt.Errorf("failed to parse schema response: %w", err)
		}

		// Create Avro codec
		codec, err := goavro.NewCodec(schemaResp.Schema)
		if err != nil {
			return nil, fmt.Errorf("failed to create codec: %w", err)
		}

		// Encode message to Avro binary
		avroBinary, err := codec.BinaryFromNative(nil, params.Message)
		if err != nil {
			return nil, fmt.Errorf("failed to encode to Avro: %w", err)
		}

		// Prepend magic byte (0) and schema ID (4 bytes, big-endian)
		messageBytes = make([]byte, 5+len(avroBinary))
		messageBytes[0] = 0 // Magic byte
		binary.BigEndian.PutUint32(messageBytes[1:5], uint32(schemaResp.ID))
		copy(messageBytes[5:], avroBinary)
	} else {
		// Plain JSON
		messageBytes = messageJSON
	}

	// Publish to Kafka
	deliveryChan := make(chan kafka.Event)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &params.Topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}, deliveryChan)

	if err != nil {
		return nil, fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return nil, fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	return &PublishResult{
		Status:    "success",
		Partition: m.TopicPartition.Partition,
		Offset:    int64(m.TopicPartition.Offset),
	}, nil
}
