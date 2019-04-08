package controllers

import (
	"encoding/json"
	"sam-book-sample/models"
	"sam-book-sample/utils"

	"github.com/aws/aws-lambda-go/events"
)

var ValidateMicropostSettings = []*ValidatorSetting{
	{ArgName: "content", ValidateTags: "required,max=140"},
}

type RequestMicropost struct {
	Content string `json:"content"`
}

type RequestPostMicropost struct {
	RequestMicropost
}

type RequestPutMicropost struct {
	RequestMicropost
}

type ResponseMicropost struct {
	ID      uint64 `json:"id"`
	UserID  uint64 `json:"user_id"`
	Content string `json:"content"`
}

type ResponseMicroposts struct {
	Microposts []*ResponseMicropost `json:"microposts"`
}

func PostMicroposts(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	validErr := ValidateBody(request.Body, ValidateMicropostSettings)
	if validErr != nil {
		return Response400(*validErr)
	}

	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	var req RequestPostMicropost
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	micropost := &models.Micropost{
		UserID:  userID,
		Content: req.Content,
	}

	err = micropost.Create()
	if err != nil {
		return Response500(err)
	}

	return Response201(micropost.ID)
}

func PutMicropost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	validErr := ValidateBody(request.Body, ValidateMicropostSettings)
	if validErr != nil {
		return Response400(*validErr)
	}

	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	var req RequestPutMicropost
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	id, err := utils.ParseUint(request.PathParameters["micropost_id"])
	if err != nil {
		return Response500(err)
	}

	micropost, err := models.GetMicropostByID(id)
	if err != nil {
		return Response500(err)
	}
	if micropost == nil {
		return Response404()
	}

	micropost.UserID = userID
	micropost.Content = req.Content

	err = micropost.Update()
	if err != nil {
		return Response500(err)
	}

	return Response200OK()
}

func GetMicroposts(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	microposts, err := models.GetMicropostsByUserID(userID)
	if err != nil {
		return Response500(err)
	}

	var resMicroposts = make([]*ResponseMicropost, len(microposts))
	for i, m := range microposts {
		resMicroposts[i] = &ResponseMicropost{
			ID:      m.ID,
			UserID:  m.UserID,
			Content: m.Content,
		}
	}

	return Response200(&ResponseMicroposts{
		Microposts: resMicroposts,
	})
}

func GetMicropost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	micropostID, err := utils.ParseUint(request.PathParameters["micropost_id"])
	if err != nil {
		return Response500(err)
	}

	micropost, err := models.GetMicropostByID(micropostID)
	if err != nil {
		return Response500(err)
	}
	if micropost == nil {
		return Response404()
	}

	res := &ResponseMicropost{
		ID:      micropost.ID,
		Content: micropost.Content,
		UserID:  micropost.UserID,
	}

	return Response200(res)
}

func DeleteMicropost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	id, err := utils.ParseUint(request.PathParameters["micropost_id"])
	if err != nil {
		return Response500(err)
	}

	err = models.DeleteMicropost(id)
	if err != nil {
		return Response500(err)
	}

	return Response200OK()
}
