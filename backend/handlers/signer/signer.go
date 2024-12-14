package signer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

// customCredentialsProvider implements aws.CredentialsProvider.
type customCredentialsProvider struct {
	AccessKeyID     string
	SecretAccessKey string
}

// Retrieve satisfies the aws.CredentialsProvider interface.
func (p customCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     p.AccessKeyID,
		SecretAccessKey: p.SecretAccessKey,
		Source:          "CustomEnvProvider",
	}, nil
}

// CreateSignedRequest generates an AWS Signature v4 signed request.
func CreateSignedRequest(method, url string, region string) (*http.Request, error) {
	// Load credentials from environment variables
	accessKey := os.Getenv("SCW_ACCESS_KEY")
	secretKey := os.Getenv("SCW_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("missing SCW_ACCESS_KEY or SCW_SECRET_KEY environment variables")
	}

	// Create a custom credentials provider
	credsProvider := customCredentialsProvider{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
	}

	// Retrieve credentials
	creds, err := credsProvider.Retrieve(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credentials: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Compute a payload hash for an empty body (GET request)
	payloadHash := sha256.Sum256([]byte{})
	payloadHashHex := fmt.Sprintf("%x", payloadHash)

	// Add the X-Amz-Content-Sha256 header
	req.Header.Set("X-Amz-Content-Sha256", payloadHashHex)


	// Sign the request using AWS Signature v4
	signer := v4.NewSigner()
	err = signer.SignHTTP(
		context.TODO(),        // Context
		creds,                 // Credentials
		req,                   // Request to sign
		payloadHashHex,        // SHA256 payload hash
		"s3",                  // Service
		region,                // Region
		time.Now(),            // Signing time
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %v", err)
	}

	return req, nil
}