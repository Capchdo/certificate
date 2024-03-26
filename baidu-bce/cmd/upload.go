package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"main/util"

	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload DOMAIN CERTIFICATE_DIRECTORY",
	Short: "上传 SSL 证书",
	Long: `上传 CERTIFICATE_DIRECTORY 中的 SSL 证书 fullchain.pem 和 privkey.pem 到指定域名 DOMAIN

可以反复上传，并且 SSL 证书的 ID 不变。

等效网页：
- https://console.bce.baidu.com/iam/#/iam/cert/list → 添加证书
- https://console.bce.baidu.com/cdn/#/cdn/list
- https://console.bce.baidu.com/cdn/#/cdn/detail/https~domain=… → HTTPS 配置 → 证书选择`,
	Args:    cobra.ExactArgs(2),
	Example: "  baidu-bce upload haobit.top /etc/letsencrypt/live/haobit.top/",
	Run: func(cmd *cobra.Command, args []string) {
		main_domain := args[0]
		certificate := read_certificate(args[1])

		client := util.BuildCDNClient()
		put_certificate(client, main_domain, certificate)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func read_certificate(directory string) api.UserCertificate {
	fullchain, err := os.ReadFile(filepath.Join(directory, "fullchain.pem"))
	if err != nil {
		panic(fmt.Errorf("failed to read fullchain: %w", err))
	}
	privkey, err := os.ReadFile(filepath.Join(directory, "privkey.pem"))
	if err != nil {
		panic(fmt.Errorf("failed to read privkey: %w", err))
	}

	return api.UserCertificate{
		CertName:    "lets-encrypt-" + time.Now().Format(time.DateOnly),
		ServerData:  string(fullchain),
		PrivateData: string(privkey),
	}
}

func put_certificate(client *cdn.Client, main_domain string, certificate api.UserCertificate) {
	// Put to the main domain
	log.Printf("Uploading the certificate…\n")
	id, err := client.PutCert(main_domain, &certificate, "ON")
	if err != nil {
		log.Fatalf("Fail to put the main certificate: %+v.\n", err)
	}
	log.Printf("Got certificate ID: %v.\n", id)

	log.Printf("Listing domains…\n")
	domains, next_marker, err := client.ListDomains("")
	if err != nil {
		log.Fatalf("Fail to list domains: %v", err)
	}
	if next_marker != "" {
		// https://github.com/baidubce/bce-sdk-go/blob/86581e5eb81df460f8f2f99b20216585711d5b78/doc/CDN.md#%E5%9F%9F%E5%90%8D%E5%88%97%E8%A1%A8%E6%9F%A5%E8%AF%A2-listdomains
		// 目前百度服务器的分页限制很大，无需考虑
		log.Fatalf("Only empty next marker is supported yet: %v.", next_marker)
	}

	// Put to other domains
	for _, domain := range domains {
		if domain == main_domain {
			// Automatically set
			continue
		}

		if strings.HasSuffix(domain, main_domain) {
			log.Printf("Putting certificate to %v…\n", domain)
			err := client.SetDomainHttps(domain, &api.HTTPSConfig{
				Enabled: true,
				CertId:  id,
			})
			if err != nil {
				log.Fatalf("Fail to put certificate to %v: %+v.\n", domain, err)
			}
		} else {
			log.Printf("Ignore %v because it does not match the main domain %v.\n", domain, main_domain)
		}

	}

	// Verify
	detail, err := client.GetCert(main_domain)
	if err != nil {
		log.Fatalf("Fail to verify the certificate of %v: %+v.\n", main_domain, err)
	}
	fmt.Printf("Successfully put to %v.\n", main_domain)
	fmt.Printf("🔒 Certificate: %v (%v)\n", detail.CertName, detail.CertId)
	fmt.Printf("🎯 Domain: %v (%v)\n", detail.CommonName, detail.DNSNames)
	fmt.Printf("📅 %v – %v\n", detail.StartTime, detail.StopTime)
	fmt.Printf("🚗 Status: %v\n", detail.Status)
	fmt.Printf("📝 Created at %v, updated at %v.\n", detail.CreateTime, detail.UpdateTime)
}
