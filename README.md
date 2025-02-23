# Huan-WebWatcher
## 介绍
简单的 HTTP/HTTPS 监控系统，通过 Get 请求检查网站是否正常，并通过状态码判断状态。

## 如何配置
### 命令行参数
```text
Usage of HSWv1.exe:
  --help
  --h
          Show usage of HSBv1.exe. If this option is set, the backend service
          will not run.

  --version
  --v
          Show version of HSBv1.exe. If this option is set, the backend service
          will not run.

  --license
  --l
          Show license of HSBv1.exe. If this option is set, the backend service
          will not run.

  --report
  --r
          Show how to report questions/errors of HSBv1.exe. If this option is
          set, the backend service will not run.

  --config string
  --c string
          The location of the running configuration file of the backend service.
          The option is a string, the default value is config.yaml in the
          running directory.

  --output-config string
          The location of the reverse output after the backend service running
          configuration file is parsed. The option is a string and the default
          is config.output.yaml in the running directory.
```

根据上面的描述，我们主要使用`--config`参数，该参数表示配置文件的位置。默认值是：`config.yaml`。

当`--config`为`config.yaml`（默认值）时，`--output-config`则会默认设置为`config.output.yaml`，并将配置文件输出到此位置。
输出的配置文件是完整版，包含全部选项和默认选项的，同时过滤非法选项。

### 配置文件
配置文件是`yaml`文件，请看以下配置文件：

```yaml
mode: debug  # 运行模式（Debug/Release/Test）
log-level: debug  # 日志记录登记
log-tag: enable  # 是否输出标签日志（Debug使用）
time-zone: Local  # 时区（UTC/Local/指定时区），若指定时区不存在，会退化到Local（本地电脑时区），若仍不存在则退化到UTC
name: 001  # 服务名称（会显示在消息推送中）

watcher:
  urls:
    - name: '百度' # URL的名字，当URL比较长的时候可以设定名字来缩短显示的URL，若不设置则默认 name 为 url的host
      url: https://www.baidu.com  # 网站的URL（必须是https协议）
      deadline: 150d  # 即将过期的标准，若证书在 deadline 时间内过期，则会发出警告。例如此处设置为150d则表示证书在150天内过期则会发出警报。
    - name: '百度' # URL的名字，当URL比较长的时候可以设定名字来缩短显示的URL，若不设置则默认 name 为 url的host
      url: https://song-zh.com  # 网站的URL（必须是https协议或http协议）
      skip-tls-verify: disable  # 是否跳过证书检查（默认不跳过）
      status:  # 表示网站正常的状态码（默认为2xx）
        - '2xx'  # 支持 1xx 2xx 3xx 4xx 5xx，例如 1xx 则表示状态码在 100-199 范围内，以此类推
        - 'xxx'  # 表示任何状态码均可，只要网站请求成功即视为正常
        - '200'  # 具体的状态码
api:
  webhook: # 企业微信机器人 Webhook，可为空，关闭企业微信推送

smtp:  # 发送邮件消息推送
  address: # smtp 服务器地址，可为空，为空表示关闭smtp
  user: # smtp 用户名（邮件），可为空，为空表示关闭smtp
  password: # smtp 用户密码
  recipient:
    - xxx@wxample.com  # 接收邮件通知的用户
```

## 构建与运行
### 构建
使用`go build`指令进行编译。
```shell
$ go build github.com/SongZihuan/web-watcher/src/cmd/webwatcher/hwwv1
```

生产环境下可以使用一些编译标志来压缩目标文件大小。
```shell
$ go build -trimpath -ldflags='-s -w' github.com/SongZihuan/web-watcher/src/cmd/webwatcher/hwwv1
```

### 运行
执行编译好的可执行文件即可。具体命令行参数可参见上文。

## 协议
本软件基于 [MIT LICENSE](/LICENSE) 发布。
了解更多关于 MIT LICENSE , 请 [点击此处](https://mit-license.song-zh.com) 。
