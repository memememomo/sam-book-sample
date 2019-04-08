package settings

import (
	"encoding/base64"
	"os"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

var kmsClient *kms.KMS

func keyID() string {
	return os.Getenv("KMS_KEY_ID")
}

func getKMSClient() *kms.KMS {
	if kmsClient == nil {
		kmsClient = kms.New(
			session.Must(session.NewSession()),
			aws.NewConfig().WithRegion("ap-northeast-1"),
		)
	}
	return kmsClient
}

func DecryptByKMS(str string) (string, error) {
	client := getKMSClient()

	decodedBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", errors.WithStack(err)
	}
	res, err := client.Decrypt(&kms.DecryptInput{
		CiphertextBlob: decodedBytes,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(res.Plaintext[:]), nil
}

func EncryptByKMS(str string) (string, error) {
	client := getKMSClient()

	res, err := client.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyID()),
		Plaintext: []byte(str),
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	encodedString := base64.StdEncoding.EncodeToString(res.CiphertextBlob)

	return encodedString, nil
}
