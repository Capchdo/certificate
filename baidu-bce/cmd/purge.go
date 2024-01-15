package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"main/util"

	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/spf13/cobra"
)

var wait bool

var purgeCmd = &cobra.Command{
	Use:   "purge URL...",
	Short: "刷新 CDN 中的缓存",
	Long: `刷新 CDN 中的缓存

“/”结尾的按目录刷新，其余按文件刷新。默认 HTTPS。

等效网页：
- https://console.bce.baidu.com/cdn/#/cdn/refresh/url
- https://console.bce.baidu.com/cdn/#/cdn/refresh/path
- https://console.bce.baidu.com/cdn/#/cdn/refresh/history`,
	Example: "  baidu-bce purge status.haobit.top/ haobit.top/feed.rss",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tasks := []api.PurgeTask{}
		for _, u := range args {
			tasks = append(tasks, build_task(u))
		}

		client := util.BuildClient()

		log.Printf("Sending purge tasks…")
		id, err := client.Purge(tasks)
		if err != nil {
			log.Fatalf("Fail to send purge tasks: %+v.\n", err)
		}
		log.Printf("Got task ID: %v.\n", id)

		if wait {
			wait_until_completed(client, id)
		}
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)

	purgeCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Whether to wait until completed")
}

func build_task(url string) api.PurgeTask {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	if strings.HasSuffix(url, "/") {
		return api.PurgeTask{
			Url:  url,
			Type: "directory",
		}
	} else {
		return api.PurgeTask{
			Url:  url,
			Type: "file",
		}
	}
}

func wait_until_completed(client *cdn.Client, id api.PurgedId) {
	query := api.CStatusQueryData{
		Id: string(id),
	}

	all_completed := false
	for !all_completed {
		fmt.Println("Wait a second…")
		time.Sleep(1 * time.Second)

		status, err := client.GetPurgedStatus(&query)
		if err != nil {
			log.Printf("Fail to get purged status, ignore: %v.\n", err)
		}

		all_completed = true

		for _, detail := range status.Details {
			fmt.Printf("Task: %v (%v)\n", detail.Task.Url, detail.Task.Type)
			fmt.Printf("Progress: %v%%\n", detail.CachedDetail.Progress)

			all_completed = all_completed && detail.CachedDetail.Status == "completed"
		}
	}
}
