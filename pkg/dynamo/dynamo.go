package dynamo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Service represents a client that can talk to DynamoDB.
type Service interface {
	Health() map[string]string
	Close() error
}

type service struct {
	client *dynamodb.Client
}

var (
	region        = os.Getenv("AWS_DEFAULT_REGION")
	endpoint      = os.Getenv("DYNAMO_ENDPOINT")
	accessKeyID   = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey     = os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken  = os.Getenv("AWS_SESSION_TOKEN")
	instanceStore *service
)

// New creates a DynamoDB client configured via environment variables.
// It reuses a singleton instance when available.
func New() Service {
	if instanceStore != nil {
		return instanceStore
	}

	if region == "" {
		region = "us-east-1"
	}

	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if accessKeyID != "" && secretKey != "" {
		loadOpts = append(loadOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, sessionToken)))
	}

	if endpoint != "" {
		loadOpts = append(loadOpts, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
			}),
		))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		log.Fatalf("could not load AWS config: %v", err)
	}

	instanceStore = &service{
		client: dynamodb.NewFromConfig(cfg),
	}

	return instanceStore
}

// Health checks that DynamoDB is reachable by listing tables with a short timeout.
func (s *service) Health() map[string]string {
	stats := make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := s.client.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("dynamo down: %v", err)
		log.Printf("dynamo down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	return stats
}

// Close drops the cached client reference.
func (s *service) Close() error {
	log.Printf("Disconnected from dynamodb in region: %s", region)
	instanceStore = nil
	return nil
}
