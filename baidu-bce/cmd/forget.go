package cmd

import (
	"fmt"
	"log"
	"main/util"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/dns"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
)

var forgetRecord = &cobra.Command{
	Use:   "forget DOMAIN TYPE VALUE",
	Short: "删除 DNS 解析",
	Long: `删除指定 DOMAIN、TYPE、VALUE 的记录

为避免误操作相近记录，需要提供全部信息。

与record的位置参数相同。

等效网页：
- https://console.bce.baidu.com/dns/#/dns/manage/list
- https://console.bce.baidu.com/dns/#/dns/domain/list?zoneName=… → 选中 → 删除`,
	Args:    cobra.ExactArgs(3),
	Example: "  baidu-bce forget _acme-challenge.haobit.top TXT 6kSGMVJoOhx1YMM-xc",
	Run: func(cmd *cobra.Command, args []string) {
		parts := strings.Split(args[0], ".")
		main_domain := strings.Join(parts[len(parts)-2:], ".")
		sub_domain := strings.Join(parts[:len(parts)-2], ".")

		type_ := args[1]
		value := args[2]

		client := util.BuildDNSClient()
		record := match_record(client, main_domain, type_, sub_domain, value)

		err := client.DeleteRecord(main_domain, record.Id, "") // 目前没有必要保证幂等性，故未设置 client token
		if err != nil {
			log.Fatalf("Fail to forget the record: %+v.\n", err)
		}
		fmt.Printf("Successfully forget the %v record %v with value “%v” on %v.\n", type_, sub_domain, value, main_domain)
	},
}

func init() {
	rootCmd.AddCommand(forgetRecord)
}

func match_record(client *dns.Client, main_domain string, type_ string, sub_domain string, value string) dns.Record {
	// 目前记录不多，故未考虑分页
	result, err := client.ListRecord(main_domain, &dns.ListRecordRequest{
		Rr: sub_domain,
	})
	if err != nil {
		log.Fatalf("Fail to list records: %+v.\n", err)
	}

	matched := []dns.Record{}
	for _, r := range result.Records {
		if r.Type == type_ && r.Value == value {
			matched = append(matched, r)
		}
	}

	if len(matched) == 0 {
		log.Fatalf("Fail to match any record.\n")
	} else if len(matched) > 1 {
		pp.Print(matched)
		log.Fatalf("Expect to match a single record, but there are %v records.\n", len(matched))
	}

	return matched[0]
}
