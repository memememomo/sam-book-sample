package controllers

import (
	"encoding/json"
	"fmt"
	"sam-book-sample/db"
	"sam-book-sample/mocks"
	"sam-book-sample/models"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestPostUsers_201(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	body := map[string]interface{}{
		"user_name": "テスト名前",
		"email":     "test@example.com",
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	res := PostUsers(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
	})

	assert.Equal(t, 201, res.StatusCode)

	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	id := uint64(resBody["id"].(float64))
	assert.Equal(t, uint64(1), id)

	user, err := models.GetUserByID(id)
	assert.NoError(t, err)
	assert.Equal(t, body["user_name"].(string), user.Name)
	assert.Equal(t, body["email"].(string), user.Email)
}

func TestPostUsers_400(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	userGen := mocks.User()
	userMock := userGen.Single(0, mocks.E).(*models.User)

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		{
			Request: map[string]interface{}{
				"user_name": "",
				"email":     "",
			},
			Expected: map[string]interface{}{
				"user_name": "ユーザー名を入力してください。",
				"email":     "メールアドレスを入力してください。",
			},
		},
		{
			Request: map[string]interface{}{
				"user_name": "hoge",
				"email":     "test@",
			},
			Expected: map[string]interface{}{
				"email": "メールアドレスの形式が不正です。",
			},
		},
		{
			Request: map[string]interface{}{
				"user_name": "dup",
				"email":     userMock.Email,
			},
			Expected: map[string]interface{}{
				"email": "すでに登録されているメールアドレスです。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PostUsers(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
		})

		var resBody map[string]interface{}
		err = json.Unmarshal([]byte(res.Body), &resBody)
		assert.NoError(t, err)

		errors := resBody["errors"].(map[string]interface{})

		assert.Equal(t, 400, res.StatusCode, msg)
		assert.Equal(t, c.Expected, errors)
	}
}

func TestPutUser_200(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	userGen := mocks.User()
	userMock := userGen.Single(0, mocks.E).(*models.User)

	body := map[string]interface{}{
		"user_name": "テスト名前更新",
		"email":     "test_update@example.com",
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	res := PutUser(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	assert.Equal(t, 200, res.StatusCode)

	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	user, err := models.GetUserByID(userMock.ID)
	assert.NoError(t, err)
	assert.Equal(t, body["user_name"].(string), user.Name)
	assert.Equal(t, body["email"].(string), user.Email)
}

func TestPutUser_200_dup(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	userGen := mocks.User()
	userMock := userGen.Single(0, mocks.E).(*models.User)

	body := map[string]interface{}{
		"user_name": "テスト名前更新",
		"email":     userMock.Email,
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	res := PutUser(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	assert.Equal(t, 200, res.StatusCode)

	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	user, err := models.GetUserByID(userMock.ID)
	assert.NoError(t, err)
	assert.Equal(t, body["user_name"].(string), user.Name)
	assert.Equal(t, body["email"].(string), user.Email)
}

func TestPutUser_400(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	userGen := mocks.User()
	userMock := userGen.Single(0, mocks.E).(*models.User)

	dupUserMock := userGen.Single(1, mocks.E).(*models.User)

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		{
			Request: map[string]interface{}{
				"user_name": "",
				"email":     "",
			},
			Expected: map[string]interface{}{
				"user_name": "ユーザー名を入力してください。",
				"email":     "メールアドレスを入力してください。",
			},
		},
		{
			Request: map[string]interface{}{
				"user_name": "hoge",
				"email":     "test@",
			},
			Expected: map[string]interface{}{
				"email": "メールアドレスの形式が不正です。",
			},
		},
		{
			Request: map[string]interface{}{
				"user_name": "hoge",
				"email":     dupUserMock.Email,
			},
			Expected: map[string]interface{}{
				"email": "すでに登録されているメールアドレスです。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PutUser(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
			PathParameters: map[string]string{
				"user_id": fmt.Sprintf("%d", userMock.ID),
			},
		})

		var resBody map[string]interface{}
		err = json.Unmarshal([]byte(res.Body), &resBody)
		assert.NoError(t, err)

		errors := resBody["errors"].(map[string]interface{})

		assert.Equal(t, 400, res.StatusCode, msg)
		assert.Equal(t, c.Expected, errors)
	}
}

func TestGetUser(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	userGen := mocks.User()
	userMock := userGen.Single(0, mocks.E).(*models.User)

	res := GetUser(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	assert.Equal(t, 200, res.StatusCode)

	var body map[string]interface{}
	err := json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)
	assert.Equal(t, float64(userMock.ID), body["id"])
	assert.Equal(t, userMock.Name, body["user_name"])
	assert.Equal(t, userMock.Email, body["email"])
}

func TestGetUsers(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	userGen := mocks.User()
	userMocks := userGen.Multi(2, mocks.E)

	res := GetUsers(events.APIGatewayProxyRequest{})

	assert.Equal(t, 200, res.StatusCode)

	var body map[string]interface{}
	err := json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)

	actualUsers := body["users"].([]interface{})

	expected1 := userMocks[1].(*models.User)
	actual1 := actualUsers[0].(map[string]interface{})
	assert.Equal(t, float64(expected1.ID), actual1["id"])
	assert.Equal(t, expected1.Name, actual1["user_name"])
	assert.Equal(t, expected1.Email, actual1["email"])

	expected2 := userMocks[0].(*models.User)
	actual2 := actualUsers[1].(map[string]interface{})
	assert.Equal(t, float64(expected2.ID), actual2["id"])
	assert.Equal(t, expected2.Name, actual2["user_name"])
	assert.Equal(t, expected2.Email, actual2["email"])
}

func TestDeleteUser(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	userGen := mocks.User()
	userMock := userGen.Single(0, mocks.E).(*models.User)

	res := DeleteUser(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userMock.ID),
		},
	})

	assert.Equal(t, 200, res.StatusCode)

	users, err := models.GetUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}
