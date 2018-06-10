package groupbot

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

var intentRouter_ = IntentRouter{}

func ServeAwsLambda(intentRouter IntentRouter) {
	log.Info("serveAwsLambda_S1")
	intentRouter_ = intentRouter
	lambda.Start(HandleAwsLambda)
}

func HandleAwsLambda(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Info("HandleAwsLambda_S1")
	log.Info(fmt.Sprintf("HandleAwsLambda_S2_req_body: %v", req.Body))

	if val, ok := req.Headers[ValidationTokenHeader]; ok && len(val) > 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{ValidationTokenHeader: val},
			Body:       `{"statusCode":200}`,
		}, nil
	}

	bot := Groupbot{IntentRouter: intentRouter_}
	return bot.HandleAwsLambda(req)
}
