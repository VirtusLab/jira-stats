package analyzer

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/ztrue/tracerr"
	"log"
)

func RetrieveSecrets() ([]byte, error) {
	secretId := "JiraCreds"

	sess := session.Must(session.NewSession())
	secretMgr := secretsmanager.New(sess)

	output, err := secretMgr.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretId})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	log.Printf("Fetched Secret id: %s\n", *output.Name)

	return []byte(*output.SecretString), nil
}
