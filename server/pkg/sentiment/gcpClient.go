package sentiment

import (
	"context"
	"fmt"
	"log"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

type GcpClient struct {
	languageClient *language.Client
}

func NewGcpClient() *GcpClient {
	languageClient, err := language.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create languageClient: %v", err)
	}
	return &GcpClient{
		languageClient: languageClient,
	}
}

func (client *GcpClient) CloseClient() {
	client.languageClient.Close()
}

func (client *GcpClient) CalculateSentiment(
	text string,
) (float32, error) {
	// Detects the sentiment of the text.
	sentiment, err := client.languageClient.AnalyzeSentiment(context.Background(), &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		return 0, fmt.Errorf("Failed to analyze text: %v", err)
	}
	return sentiment.DocumentSentiment.Score, nil
}
