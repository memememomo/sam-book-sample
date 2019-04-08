package db

import (
	"fmt"
	"os"
	"sam-book-sample/settings"
	"sam-book-sample/utils"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/golang/glog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

const PKName = "PK"
const SKName = "SK"

type MainTable struct {
	PK string `dynamo:"PK,hash"`
	SK string `dynamo:"SK,range"`
}

func Config() *aws.Config {
	config := &aws.Config{
		Region: aws.String("ap-northeast-1"),
	}

	envs := settings.Env()

	if envs.IsDynamoLocal() {
		config.Credentials = credentials.NewStaticCredentials("dummy", "dummy", "dummy")
		config.Endpoint = aws.String(envs.DynamoEndpoint())
	}

	return config
}

func ConnectDB() (*dynamo.DB, error) {
	dynamoSession, err := session.NewSession(Config())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return dynamo.New(dynamoSession), nil
}

func Table() (*dynamo.Table, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	table := db.Table(settings.Env().DynamoTableName())
	return &table, nil
}

func DropTable() error {
	db, err := ConnectDB()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.Client().
		DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(settings.Env().DynamoTableName()),
		})

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func MigrateForTest() error {
	db, err := ConnectDB()
	if err != nil {
		return errors.WithStack(err)
	}

	err = db.
		CreateTable(settings.Env().DynamoTableName(), MainTable{}).
		Provision(100, 100).
		Run()

	if err != nil {
		return errors.WithStack(err)
	}

	if settings.Env().IsDebug() {
		err = DescribeTable()
		if err != nil {
			glog.Warningf("Failed to describe table: %s", err.Error())
		}
	}

	return nil
}

func DescribeTable() error {
	table, err := Table()
	if err != nil {
		return errors.WithStack(err)
	}

	desc, err := table.Describe().Run()
	if err != nil {
		return errors.WithStack(err)
	}

	glog.Infof("%#v", desc)

	return nil
}

func SetupDBForTest() error {
	// テスト毎に新しいテーブルを作成するため、ランダムな値を設定する
	randomStr, err := utils.GenerateToken(30)
	if err != nil {
		return errors.WithStack(err)
	}
	os.Setenv("DYNAMO_TABLE_VERSION", randomStr)
	return MigrateForTest()
}

func GenerateID(tableName string) (uint64, error) {
	attr, err := AtomicCount(fmt.Sprintf("AtomicCounter-%s", tableName), "AtomicCounter", "CurrentNumber", 1)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	nStr := aws.StringValue(attr.N)
	n, err := utils.ParseUint(nStr)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return n, nil
}

func AtomicCount(pk, sk, counterName string, value int) (*dynamodb.AttributeValue, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	output, err := db.Client().UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(settings.Env().DynamoTableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			counterName: {
				Action: aws.String("ADD"),
				Value: &dynamodb.AttributeValue{
					N: aws.String(fmt.Sprintf("%d", value)),
				},
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return output.Attributes[counterName], nil
}
