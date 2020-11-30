package secrets

import (
	"fmt"
	"os"
)

func GetSecret(secretID string) (string, error) {
	apiKey, exists := os.LookupEnv(secretID)

	if exists {
		return apiKey, nil
	}

	return "", fmt.Errorf("%s was not set", secretID)
}
