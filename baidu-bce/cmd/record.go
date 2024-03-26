package cmd

import (
	"fmt"
	"log"
	"main/util"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/dns"
	"github.com/spf13/cobra"
)

var description string

var recordCmd = &cobra.Command{
	Use:   "record DOMAIN TYPE VALUE",
	Short: "添加 DNS 解析",
	Long: `向指定域名 DOMAIN 记录指定类型 TYPE 的 VALUE

可以记录相同域名、类型的不同值；但若值也相同，则不能记录。

另外，百度的API禁止相同域名同时有CNAME、TXT两种类型的记录，但网页上允许。

与record的位置参数相同。

等效网页：
- https://console.bce.baidu.com/dns/#/dns/manage/list
- https://console.bce.baidu.com/dns/#/dns/domain/list?zoneName=… → 添加解析`,
	Args:    cobra.ExactArgs(3),
	Example: "  baidu-bce record _acme-challenge.haobit.top TXT 6kSGMVJoOhx1YMM-xc",
	Run: func(cmd *cobra.Command, args []string) {
		parts := strings.Split(args[0], ".")
		main_domain := strings.Join(parts[len(parts)-2:], ".")
		sub_domain := strings.Join(parts[:len(parts)-2], ".")

		type_ := args[1]
		value := args[2]

		client := util.BuildDNSClient()
		err := client.CreateRecord(main_domain, &dns.CreateRecordRequest{
			Type:        type_,
			Rr:          sub_domain,
			Value:       value,
			Description: &description,
		}, "") // 目前没有必要保证幂等性，故未设置 client token
		if err != nil {
			log.Fatalf("Fail to record: %+v.\n", err)
		}
		fmt.Printf("Successfully forget the %v record %v with value “%v” on %v.\n", type_, sub_domain, value, main_domain)
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.Flags().StringVar(&description, "description", "", "description of the record (no default)")
}
