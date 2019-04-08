package controllers

import (
	"encoding/json"
	"fmt"
	"sam-book-sample/db"
	"sam-book-sample/mocks"
	"sam-book-sample/models"
	"strings"
	"testing"

	"github.com/memememomo/dbmock"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
)

func TestPostMicroposts_201(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	body := map[string]interface{}{
		"content": strings.Repeat("a", 140),
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	userID := uint64(1)

	res := PostMicroposts(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	})

	assert.Equal(t, 201, res.StatusCode)

	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	id := uint64(resBody["id"].(float64))
	assert.Equal(t, uint64(1), id)

	micropost, err := models.GetMicropostByID(id)
	assert.NoError(t, err)
	assert.Equal(t, body["content"].(string), micropost.Content)
	assert.Equal(t, userID, micropost.UserID)
}

func TestPostMicroposts_400(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		{
			Request: map[string]interface{}{
				"content": "",
			},
			Expected: map[string]interface{}{
				"content": "本文を入力してください。",
			},
		},
		{
			Request: map[string]interface{}{
				"content": strings.Repeat("a", 141),
			},
			Expected: map[string]interface{}{
				"content": "本文の文字数が上限を超えています。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PostMicroposts(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
			PathParameters: map[string]string{
				"user_id": "1",
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

func TestPutMicropost_200(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	micropostGen := mocks.Micropost()
	micropostMock := micropostGen.Single(0, mocks.E).(*models.Micropost)

	body := map[string]interface{}{
		"content": strings.Repeat("a", 140),
	}
	bodyStr, err := json.Marshal(body)
	assert.NoError(t, err)

	res := PutMicropost(events.APIGatewayProxyRequest{
		Body: string(bodyStr),
		PathParameters: map[string]string{
			"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
			"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
		},
	})

	assert.Equal(t, 200, res.StatusCode)

	var resBody map[string]interface{}
	err = json.Unmarshal([]byte(res.Body), &resBody)
	assert.NoError(t, err)

	micropost, err := models.GetMicropostByID(micropostMock.ID)
	assert.NoError(t, err)
	assert.Equal(t, body["content"].(string), micropost.Content)
}

func TestPutMicropost_400(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	micropostGen := mocks.Micropost()
	micropostMock := micropostGen.Single(0, mocks.E).(*models.Micropost)

	cases := []struct {
		Request  map[string]interface{}
		Expected map[string]interface{}
	}{
		{
			Request: map[string]interface{}{
				"content": "",
			},
			Expected: map[string]interface{}{
				"content": "本文を入力してください。",
			},
		},
		{
			Request: map[string]interface{}{
				"content": strings.Repeat("a", 141),
			},
			Expected: map[string]interface{}{
				"content": "本文の文字数が上限を超えています。",
			},
		},
	}

	for i, c := range cases {
		msg := fmt.Sprintf("Case:%d", i+1)

		body := c.Request
		bodyStr, err := json.Marshal(body)
		assert.NoError(t, err)

		res := PutMicropost(events.APIGatewayProxyRequest{
			Body: string(bodyStr),
			PathParameters: map[string]string{
				"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
				"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
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

func TestGetMicropost(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	micropostGen := mocks.Micropost()
	micropostMock := micropostGen.Single(0, mocks.E).(*models.Micropost)

	res := GetMicropost(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
			"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
		},
	})

	var body map[string]interface{}
	err := json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)
	assert.Equal(t, float64(micropostMock.ID), body["id"])
	assert.Equal(t, micropostMock.Content, body["content"])
	assert.Equal(t, float64(micropostMock.UserID), body["user_id"])
}

func TestGetMicroposts(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	micropostGen := mocks.Micropost()
	micropostMocks := micropostGen.Multi(3, func(i uint64, mapper dbmock.DBMapper) dbmock.DBMapper {
		m := mapper.(*models.Micropost)
		if i == 2 {
			m.UserID = 2
		} else {
			m.UserID = 1
		}
		return m
	})

	res := GetMicroposts(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id": "1",
		},
	})

	assert.Equal(t, 200, res.StatusCode)

	var body map[string]interface{}
	err := json.Unmarshal([]byte(res.Body), &body)
	assert.NoError(t, err)

	actualMicroposts := body["microposts"].([]interface{})

	expected1 := micropostMocks[0].(*models.Micropost)
	actual1 := actualMicroposts[0].(map[string]interface{})
	assert.Equal(t, float64(expected1.ID), actual1["id"])
	assert.Equal(t, expected1.Content, actual1["content"])
	assert.Equal(t, float64(expected1.UserID), actual1["user_id"])

	expected2 := micropostMocks[0].(*models.Micropost)
	actual2 := actualMicroposts[0].(map[string]interface{})
	assert.Equal(t, float64(expected2.ID), actual2["id"])
	assert.Equal(t, expected2.Content, actual2["content"])
	assert.Equal(t, float64(expected2.UserID), actual2["user_id"])
}

func TestDeleteMicropost(t *testing.T) {
	mocks.SetupDB(t)
	defer db.DropTable()

	micropostGen := mocks.Micropost()
	micropostMock := micropostGen.Single(0, mocks.E).(*models.Micropost)

	res := DeleteMicropost(events.APIGatewayProxyRequest{
		PathParameters: map[string]string{
			"user_id":      fmt.Sprintf("%d", micropostMock.UserID),
			"micropost_id": fmt.Sprintf("%d", micropostMock.ID),
		},
	})

	assert.Equal(t, 200, res.StatusCode)

	microposts, err := models.GetMicropostsByUserID(micropostMock.UserID)
	assert.NoError(t, err)
	assert.Len(t, microposts, 0)
}
