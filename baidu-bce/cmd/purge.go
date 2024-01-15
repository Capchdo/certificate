package cmd

import (
	"log"
	"strings"
	"time"

	"main/util"

	"github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
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
	Example: "  baidu-bce purge --wait status.haobit.top/ haobit.top/feed.rss",
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

	// 1. Initialize progress bars

	// Get status for the first time to determine number of progress bars
	status, err := client.GetPurgedStatus(&query)
	if err != nil {
		log.Fatalf("Fail to get purged status: %v.\n", err)
	}

	// To support color in Windows following both options are required
	progress := mpb.New(
		mpb.WithOutput(color.Output),
		mpb.WithAutoRefresh(),
	)
	bars := []*mpb.Bar{}
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)

	for _, task := range status.Details {
		b := progress.AddBar(
			100,
			mpb.PrependDecorators(
				decor.Name(task.Task.Url, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
				decor.OnCompleteMeta(
					decor.OnComplete(
						decor.Meta(decor.Name("Purging", decor.WCSyncSpaceR), to_meta_fn(red)),
						"Completed",
					),
					to_meta_fn(green),
				),
			),
			mpb.BarFillerClearOnComplete(),
			mpb.AppendDecorators(
				decor.OnComplete(decor.Percentage(), ""),
			),
		)
		bars = append(bars, b)
	}

	// 2. Wait

	for {
		time.Sleep(1 * time.Second)

		status, err := client.GetPurgedStatus(&query)
		if err != nil {
			log.Printf("Fail to get purged status, ignore: %v.\n", err)
		}

		all_completed := true
		for i, task := range status.Details {
			// 顺序总固定
			bars[i].SetCurrent(task.CachedDetail.Progress)

			all_completed = all_completed && task.CachedDetail.Status == "completed"
		}
		if all_completed {
			break
		}
	}

	progress.Wait()
}

func to_meta_fn(c *color.Color) func(string) string {
	return func(s string) string {
		return c.Sprint(s)
	}
}
