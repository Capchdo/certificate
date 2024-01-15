package util

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/spf13/viper"
)

func BuildClient() *cdn.Client {
	access_key := viper.GetString("access_key")
	secret_key := viper.GetString("secret_key")
	client, err := cdn.NewClient(access_key, secret_key, "https://cdn.baidubce.com")
	if err != nil {
		panic(fmt.Errorf("fail to build client: %w", err))
	}

	return client
}
