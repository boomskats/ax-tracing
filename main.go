package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	token := strings.TrimPrefix(event.Headers["authorization"], "Bearer ")

    // if token is blank check the uppercased Authorization header 
    if token == "" {
        token = strings.TrimPrefix(event.Headers["Authorization"], "Bearer ")
    }

    // if token is still blank return unauthorized
    if token == "" {
        return events.APIGatewayV2CustomAuthorizerSimpleResponse{
            IsAuthorized: false,
        }, nil
    }

	userID, err := ValidateToken(token)
	if err != nil {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context: map[string]interface{}{
			"userId": userID,
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
