package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := strings.TrimPrefix(event.AuthorizationToken, "Bearer ")

	userID, err := ValidateToken(token)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	return generatePolicy(userID, "Allow", event.MethodArn, userID), nil
}

func generatePolicy(principalID, effect, resource, userId string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

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

	authResponse.Context = map[string]interface{}{
		"userId": userId,
	}

	return authResponse
}

func main() {
	lambda.Start(handleRequest)
}
