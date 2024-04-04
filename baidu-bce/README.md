# [百度智能云][bce]命令行界面

使用SSL证书、DNS（“域名服务 BCD”）、CDN。

[bce]: https://console.bce.baidu.com "百度智能云控制台"

```shell
$ baidu-bce --help
百度智能云的命令行界面

使用SSL证书、DNS（“域名服务 BCD”）、CDN。

等效网页：https://console.bce.baidu.com

Usage:
  baidu-bce [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  forget      删除 DNS 解析
  help        Help about any command
  purge       刷新 CDN 中的缓存
  record      添加 DNS 解析
  upload      上传 SSL 证书

Flags:
      --config string   config file containing access_key and secret_key (default is ./baidu-bce.yaml)
  -h, --help            help for baidu-bce

Use "baidu-bce [command] --help" for more information about a command.
```
