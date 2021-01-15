package store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
	"reflect"
)

type AWS struct {
	store *ssm.SSM
}

func NewAWS() (*AWS, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &AWS{ssm.New(sess)}, nil
}

func (a *AWS) GetParam(param string, secret bool) (string, error) {
	if secret {
		param = "/aws/reference/secretsmanager/" + param
	}

	result, err := a.store.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(param),
		WithDecryption: aws.Bool(secret),
	})
	if err != nil {
		log.Println(fmt.Sprintf("%+v", reflect.TypeOf(err)))
		return "", err
	}

	if result.Parameter.Value != nil {
		return *result.Parameter.Value, nil
	}

	return "", nil
}

func (a *AWS) PutParam(param string, value string) error {
	_, err := a.store.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(param),
		Value:     aws.String(value),
		Overwrite: aws.Bool(true),
	})

	return err
}
