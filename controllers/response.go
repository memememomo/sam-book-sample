package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/glog"
)

type Response201Body struct {
	Message string `json:"message"`
	ID      uint64 `json:"id"`
}

type Response400Body struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type Response401Body struct {
	Message string `json:"message"`
}

func commonHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
}

func Response200(body interface{}) events.APIGatewayProxyResponse {
	b, err := json.Marshal(body)
	if err != nil {
		return Response500(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
		Headers:    commonHeaders(),
	}
}

func Response200OK() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    commonHeaders(),
		Body:       `{"message":"OK"}`,
	}
}

func Response201(id uint64) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    commonHeaders(),
		Body:       fmt.Sprintf(`{"message":"OK","id":%d}`, id),
	}
}

func Response400(errs map[string]error) events.APIGatewayProxyResponse {
	glog.Warningf("%+v", errs)
	res := &Response400Body{
		Message: "入力値を確認してください。",
		Errors:  ConvertErrorsToMessage(errs),
	}

	b, err := json.Marshal(res)
	if err != nil {
		return Response500(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers:    commonHeaders(),
		Body:       string(b),
	}
}

func Response404() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers:    commonHeaders(),
		Body:       `{"message":"結果が見つかりません。"}`,
	}
}

func Response500(err error) events.APIGatewayProxyResponse {
	glog.Errorf("%+v\n", err)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers:    commonHeaders(),
		Body:       `{"message":"サーバエラーが発生しました。"}`,
	}
}
