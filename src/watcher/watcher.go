package watcher

import (
	"crypto/tls"
	"fmt"
	"github.com/SongZihuan/web-watcher/src/config"
	"github.com/SongZihuan/web-watcher/src/logger"
	"github.com/SongZihuan/web-watcher/src/notify"
	"net/http"
)

func Run() error {
	if !config.IsReady() {
		panic("config is not ready")
	}

MainCycle:
	for _, url := range config.GetConfig().Watcher.URLs {
		logger.Infof("开始请求 %s", url.Name)

		statusCode, err := httpGet(url.URL, url.SkipTLSVerify.IsEnable(false))
		if err != nil {
			logger.Errorf("请求 %s 出现异常：%s", url.Name, err.Error())
			notify.NewRecord(url.Name, url.URL, err.Error())
			continue MainCycle
		}

		for _, s := range url.Status {
			switch s {
			case "xxx":
				logger.Infof("处理 %s 完成", url.Name)
				continue MainCycle
			case "1xx":
				if statusCode >= 100 && statusCode <= 199 {
					logger.Infof("处理 %s 完成", url.Name)
					continue MainCycle
				}
			case "2xx":
				if statusCode >= 200 && statusCode <= 299 {
					logger.Infof("处理 %s 完成", url.Name)
					continue MainCycle
				}
			case "3xx":
				if statusCode >= 300 && statusCode <= 399 {
					logger.Infof("处理 %s 完成", url.Name)
					continue MainCycle
				}
			case "4xx":
				if statusCode >= 400 && statusCode <= 499 {
					logger.Infof("处理 %s 完成", url.Name)
					continue MainCycle
				}
			case "5xx":
				if statusCode >= 500 && statusCode <= 599 {
					logger.Infof("处理 %s 完成", url.Name)
					continue MainCycle
				}
			default:
				status := fmt.Sprintf("%d", statusCode)
				if status == s {
					logger.Infof("处理 %s 完成", url.Name)
					continue MainCycle
				}
			}
		}

		errMsg := fmt.Sprintf("错误的状态码 （%d）", statusCode)
		logger.Errorf("请求 %s 出现异常：%s", url.Name, errMsg)
		notify.NewRecord(url.Name, url.URL, errMsg)

		logger.Infof("处理 %s 完成", url.Name)
	}

	notify.SendNotify()

	return nil
}

func httpGet(url string, skipTLSVerify bool) (int, error) {
	// 创建一个自定义的Transport，这样我们可以访问TLS连接状态
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify}, // 忽略服务器证书验证
	}

	// 使用自定义的Transport创建一个HTTP客户端
	client := &http.Client{Transport: tr}

	// 发送请求
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("发送 Get 请求错误：%s", err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp.StatusCode, nil
}
