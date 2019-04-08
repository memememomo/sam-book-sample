package main

import (
	"sam-book-sample/controllers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return controllers.GetMicroposts(request), nil
}

func main() {
	lambda.Start(handler)
}
