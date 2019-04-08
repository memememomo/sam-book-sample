package settings

import (
	"fmt"
	"os"

	"github.com/golang/glog"

	"github.com/aws/aws-sdk-go/service/kms"
)

type Envs struct {
	KMSClient *kms.KMS
	Cache     map[string]string
}

var envs *Envs

func NewEnv() *Envs {
	return &Envs{
		Cache: make(map[string]string),
	}
}

func Env() *Envs {
	if envs == nil {
		envs = NewEnv()
	}
	return envs
}

func (c *Envs) decrypt(key string) string {
	if os.Getenv("DISABLE_ENV_DECRYPT") != "" {
		return c.env(key)
	}

	v := c.Cache[key]
	if v != "" {
		return v
	}

	str := os.Getenv(key)
	if str == "" {
		return ""
	}

	v, err := DecryptByKMS(str)
	if err != nil {
		glog.Warning(err.Error())
		return ""
	}

	c.Cache[key] = v

	return c.Cache[key]
}

func (c *Envs) env(key string) string {
	return os.Getenv(key)
}

func (c *Envs) ProjectName() string {
	return c.env("PROJECT_NAME")
}

func (c *Envs) IsDebug() bool {
	return c.env("DEBUG") != ""
}

func (c *Envs) DynamoEndpoint() string {
	return c.env("DYNAMO_ENDPOINT")
}

func (c *Envs) DynamoTableName() string {
	if c.env("TEST_DYNAMO_TABLE_NAME") != "" {
		return c.env("TEST_DYNAMO_TABLE_NAME")
	}
	return fmt.Sprintf("%s-%s-%s", c.ProjectName(), c.env("DYNAMO_TABLE_NAME"), c.env("DYNAMO_TABLE_VERSION"))
}

func (c *Envs) IsDynamoLocal() bool {
	return c.DynamoEndpoint() != ""
}
