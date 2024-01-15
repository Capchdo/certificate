# Certificate

![Updated](https://img.shields.io/endpoint?url=https%3A%2F%2Fstatus.haobit.top%2Fupload-cert)

申请、上传证书。

## 申请证书 🆕

使用 [certbot](https://eff-certbot.readthedocs.io/en/stable/using.html) 申请证书。

### 安装

First you should install [certbot](https://certbot.eff.org/instructions?ws=nginx&os=ubuntufocal).

```shell
$ sudo snap install --classic certbot
$ sudo ln -s /snap/bin/certbot /usr/bin/certbot
```

### 新增证书

调用 certbot 申请证书。

```shell
# 必须以主域名开头
$ sudo just get-cert "haobit.top,*.haobit.top,*.app.haobit.top"
$ sudo just get-cert "capchdo.com,capchdo.cn,*.capchdo.com,*.capchdo.cn"
```

然后按照提示到“[智能云解析 - 百度智能云控制台](https://console.bce.baidu.com/dns/#/dns/domain/list?zoneName=haobit.top)”添加 TXT 记录解析，到`/path/to/acme-challenge`设置 HTTP 验证。

### 检查证书

查看证书名、域名、有效期等。

```shell
$ sudo just show-cert
```

## 上传证书 ☁️

用`baidu-bce`部署到服务器。

### 安装

```shell
$ cd baidu-bce
$ go build -o baidu-bce.exe
$ $env:GOOS = 'linux'; go build -o baidu-bce && scp baidu-bce …; $env:GOOS = ''
```

### 准备文件

- `/path/to/certificate/{fullchain,privkey}.pem`

  certbot 申请得到的证书。（可`ln -s /PATH/TO/CERTBOT/DIR/live/CERT_NAME cert/NAME`）

- `./baidu-bce.yaml`

  百度智能云的`access_key`和`secret_key`。

### 上传

```shell
$ sudo ./baidu-bce upload haobit.top /etc/letsencrypt/live/haobit.top/
```

## 自动化 🚀

### 申请证书

安装 certbot 时已经设置了 cron job 或 systemd timer 自动更新。（之后不用再运行 certbot）

然而尚未设置 manual hook，不能用。

> **Note**
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
