package cmd

import (
	"fmt"
	"log"
	"main/util"

	"github.com/baidubce/bce-sdk-go/services/dns"
	"github.com/spf13/cobra"
)

var description string

var recordCmd = &cobra.Command{
	Use:   "record DOMAIN TYPE RESOURCE_RECORD VALUE",
	Short: "添加 DNS 解析",
	Long: `向指定域名 DOMAIN 的 RESOURCE_RECORD 记录指定类型 TYPE 的 VALUE

可以记录相同域名、RR（resource record）、类型的不同值；但若值也相同，则不能记录。

另外，百度的API禁止相同RR同时有CNAME、TXT两种类型的记录，但网页上允许。

等效网页：
- https://console.bce.baidu.com/dns/#/dns/manage/list
- https://console.bce.baidu.com/dns/#/dns/domain/list?zoneName=… → 添加解析`,
	Args:    cobra.ExactArgs(4),
	Example: "  baidu-bce record haobit.top TXT _acme-challenge 6kSGMVJoOhx1YMM-xc",
	Run: func(cmd *cobra.Command, args []string) {
		main_domain := args[0]
		sub_domain := args[2]
		request := &dns.CreateRecordRequest{
			Type:        args[1],
			Rr:          sub_domain,
			Value:       args[3],
			Description: &description,
		}

		client := util.BuildDNSClient()
		err := client.CreateRecord(main_domain, request, "") // 目前没有必要保证幂等性，故未设置 client token
		if err != nil {
			log.Fatalf("Fail to record: %+v.\n", err)
		}
		fmt.Printf("Successfully record %v for %v.\n", sub_domain, main_domain)
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.Flags().StringVar(&description, "description", "", "description of the record (no default)")
}
