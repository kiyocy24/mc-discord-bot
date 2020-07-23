package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// Publish
func Publish(projectID, topicID string, msg string) (string, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	t.PublishSettings.NumGoroutines = 1

	result := t.Publish(ctx, &pubsub.Message{Data: []byte(msg)})
	id, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("Get: %v", err)
	}

	return id, nil
}
