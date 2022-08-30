package ssmparam

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func (kv *Store) MustGet(key string) string {
	v, err := kv.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}

func (kv *Store) Get(key string) (string, error) {
	ssmPath := os.Getenv("SSM_PATH")
	if ssmPath == "" {
		return "", errors.New("SSM_PATH not set")
	}
	p := path.Join(ssmPath, key)

	req := ssm.GetParameterInput{
		Name:           &p,
		WithDecryption: aws.Bool(true),
	}

	resp, err := kv.client.GetParameter(&req)
	if err != nil {
		return "", fmt.Errorf("read key %s err: %w", key, err)
	}
	val := resp.Parameter.Value
	if val == nil {
		return "", errors.New("value is nil")
	}
	return *val, nil
}

func New(client *ssm.SSM) *Store {
	return &Store{
		client: client,
	}
}

type Store struct {
	client *ssm.SSM
}
