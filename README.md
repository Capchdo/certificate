# Certificate

[![Build](https://github.com/Capchdo/certificate/actions/workflows/build.yml/badge.svg)](https://github.com/Capchdo/certificate/actions/workflows/build.yml)
![Updated](https://img.shields.io/endpoint?url=https%3A%2F%2Fstatus.haobit.top%2Fupload-cert)

利用 [certbot](https://certbot.eff.org/) 和[百度智能云](https://login.bce.baidu.com/)申请、上传证书。 

## 先决条件 🛠️

首先要[安装 certbot](https://certbot.eff.org/instructions?ws=nginx&os=ubuntufocal)。

```shell
$ sudo snap install --classic certbot
$ sudo ln -s /snap/bin/certbot /usr/bin/certbot
```

如果想要自动化，还要到百度智能云生成密钥，准备`./baidu-bce.yaml`，内容如下。

```yaml
access_key: ***
secret_key: ***
```

## 申请证书 🆕

使用 [certbot](https://eff-certbot.readthedocs.io/en/stable/using.html) 从 [Let's Encrypt](https://letsencrypt.org/zh-cn/) 申请证书。

### 新增证书

调用 certbot 申请证书。

```shell
# 必须以主域名开头
$ sudo just get-cert "haobit.top,*.haobit.top,*.app.haobit.top"
$ sudo just get-cert "capchdo.com,capchdo.cn,*.capchdo.com,*.capchdo.cn"
```

现在[`justfile`](https://just.systems/man/en/)中配置了[`certbot-hooks/*.sh`](./certbot-hooks/)，应当可以自动验证。

若未配置，需要按照提示到“[智能云解析 - 百度智能云控制台](https://console.bce.baidu.com/dns/#/dns/domain/list?zoneName=haobit.top)”添加 TXT 记录解析，到服务器上的`/path/to/acme-challenge`设置 HTTP 验证。

> [!NOTE]
>
> 一般来说，如果参数没有改变，只需要设置一次 TXT 记录，无需每次申请时修改；不过 Let's Encrypt 似乎不是这样。

### 检查证书

查看证书名、域名、有效期等。

```shell
$ sudo just show-cert
```

## 上传证书 ☁️

用`baidu-bce`上传到百度智能云，从而部署。

> [!NOTE]
>
> 在自己服务器上部署有效证书是不必须的；即使部署了，因为有CDN，用户获取到的证书也不是服务器上部署的证书。但是，仍然建议在服务器上部署有效证书或CDN供应商提供的证书，并开启CDN的有效证书验证，以保护服务器与CDN服务器间的通信。详细内容参见[配置HTTPS双向认证](https://cloud.baidu.com/doc/CDN/s/rkk7q1usv)。

### 安装

[下载可执行文件`baidu-bce`（Linux）或`baidu-bce.exe`（Windows）][latest-release]，放到任意位置并给予合适权限即可。

### 准备文件

- 之前申请到的`/path/to/certificate/{fullchain,privkey}.pem`

  certbot 申请得到的证书。

- 前述`./baidu-bce.yaml`

  百度智能云的`access_key`和`secret_key`。

### 上传

```shell
$ sudo ./baidu-bce upload haobit.top /etc/letsencrypt/live/haobit.top/
```

## 自动化 🚀

### 申请证书

安装 certbot 时已经设置了 cron job 或 systemd timer 自动更新。（之后不用再运行 certbot）

最近设置了 manual hooks，但还未完整测试过，未必能用。

> [!CAUTION]
> 
> 只是理论上如此，其实从未被验证过。

```shell
# 验证自动更新流程
$ sudo just certbot renew --dry-run

# 检查定时任务
$ systemctl list-timers
# 会有一项 snap.certbot.renew.timer。

# 查看任务状态
$ systemctl status snap.certbot.renew.timer
```

如果自动更新失败，可以直接当成[新增证书](#新增证书)重来。

### 上传证书

需要手动设置 cron job。

```shell
$ sudo crontab -e -u root
```

按[提示](https://crontab.guru/)编辑，command 一列如下。

```shell
…/baidu-bce upload … 2>&1 | logger -t cert
```

以后可从日志检查上传情况，如下。

```shell
$ sudo journalctl --system -t cert
```

## FAQ

### 为何需要 root 用户？

似乎读写日志等文件需要`sudo`。

[latest-release]: https://github.com/Capchdo/certificate/releases/latest
