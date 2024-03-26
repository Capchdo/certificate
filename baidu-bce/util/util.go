package util

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/dns"
	"github.com/spf13/viper"
)

func BuildCDNClient() *cdn.Client {
	access_key := viper.GetString("access_key")
	secret_key := viper.GetString("secret_key")
	client, err := cdn.NewClient(access_key, secret_key, "https://cdn.baidubce.com")
	if err != nil {
		panic(fmt.Errorf("fail to build CDN client: %w", err))
	}

	return client
}
func BuildDNSClient() *dns.Client {
	access_key := viper.GetString("access_key")
	secret_key := viper.GetString("secret_key")
	client, err := dns.NewClient(access_key, secret_key, "https://dns.baidubce.com")
	if err != nil {
		panic(fmt.Errorf("fail to build DNS client: %w", err))
	}

	return client
}
