package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var secret = []byte("mysecrettoken") // Replace with your actual secret token

// Help function to generate an IAM policy
func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	// Optional output with custom properties of the String, Number or Boolean type.
	authResponse.Context = map[string]interface{}{
		"stringKey":  "stringval",
		"numberKey":  123,
		"booleanKey": true,
	}
	return authResponse
}

func extractHMACSignature(e events.APIGatewayCustomAuthorizerRequestTypeRequest) string {
	return e.Headers["X-Hub-Signature-256"]
}

func isHMACSignatureValid(secret, payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	expectedSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedSig), []byte(signature))
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	signature := extractHMACSignature(event)
	if signature == "" {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("missing X-Hub-Signature-256 header")
	}

	if isHMACSignatureValid(secret, []byte(event.Body), signature) {
		return generatePolicy("user", "Allow", event.MethodArn), nil
	} else {
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
