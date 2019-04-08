package controllers

import (
	"encoding/json"
	"sam-book-sample/models"
	"sam-book-sample/utils"

	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
)

var ValidateUserSettings = []*ValidatorSetting{
	{ArgName: "user_name", ValidateTags: "required"},
	{ArgName: "email", ValidateTags: "required,email"},
}

type RequestPutUser struct {
	Name  string `json:"user_name"`
	Email string `json:"email"`
}

var ValidatorPostSetting = []*ValidatorSetting{
	{ArgName: "user_name", ValidateTags: "required"},
	{ArgName: "email", ValidateTags: "required,email"},
}

type RequestPostUser struct {
	Name  string `json:"user_name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"user_name"`
	Email string `json:"email"`
}

type UsersResponse struct {
	Users []*UserResponse `json:"users"`
}

func isUniqueEmail(email string, exceptID uint64) (bool, error) {
	user, err := models.GetUserByEmail(email)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return user == nil || user.ID == exceptID, nil
}

func PostUsers(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	validErr := ValidateBody(request.Body, ValidatorPostSetting)
	if validErr != nil {
		return Response400(*validErr)
	}

	var req RequestPostUser
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	isUniq, err := isUniqueEmail(req.Email, 0)
	if err != nil {
		return Response500(err)
	}
	if !isUniq {
		return Response400(map[string]error{
			"email": ErrUniq,
		})
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}
	err = user.Create()
	if err != nil {
		return Response500(err)
	}

	return Response201(user.ID)
}

func PutUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	validErr := ValidateBody(request.Body, ValidateUserSettings)
	if validErr != nil {
		return Response400(*validErr)
	}

	var req RequestPutUser
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response500(err)
	}

	id, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	isUniq, err := isUniqueEmail(req.Email, id)
	if err != nil {
		return Response500(err)
	}
	if !isUniq {
		return Response400(map[string]error{
			"email": ErrUniq,
		})
	}

	user, err := models.GetUserByID(id)
	if err != nil {
		return Response500(err)
	}
	if user == nil {
		return Response404()
	}

	user.Email = req.Email
	user.Name = req.Name

	err = user.Update()
	if err != nil {
		return Response500(err)
	}

	return Response200OK()
}

func GetUsers(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	users, err := models.GetUsers()
	if err != nil {
		return Response500(err)
	}

	var resUsers = make([]*UserResponse, len(users))
	for i, u := range users {
		resUsers[i] = &UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	return Response200(&UsersResponse{
		Users: resUsers,
	})
}

func GetUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	userID, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		return Response500(err)
	}
	if user == nil {
		return Response404()
	}

	res := &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return Response200(res)
}

func DeleteUser(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	id, err := utils.ParseUint(request.PathParameters["user_id"])
	if err != nil {
		return Response500(err)
	}

	err = models.DeleteUser(id)
	if err != nil {
		return Response500(err)
	}

	return Response200OK()
}
