package watcher

import (
	"crypto/tls"
	"fmt"
	"github.com/SongZihuan/web-watcher/src/config"
	"github.com/SongZihuan/web-watcher/src/logger"
	"github.com/SongZihuan/web-watcher/src/notify"
	"net/http"
	"strings"
)

func Run() error {
	if !config.IsReady() {
		panic("config is not ready")
	}

	for _, url := range config.GetConfig().Watcher.URLs {
		logger.Infof("开始请求 %s", url.Name)

		_, err := httpProcessRetry(url.URL, url.Name, url.Status, url.SkipTLSVerify.IsEnable(false))
		if err != nil {
			logger.Errorf("请求 %s 出现异常：%s", url.Name, err.Error())
			notify.NewRecord(url.Name, url.URL, err.Error())
		}

		logger.Infof("处理 %s 完成", url.Name)
	}

	notify.SendNotify()

	return nil
}

func httpProcessRetry(url string, name string, statusList []string, skipTLSVerify bool) (int, error) {
	var err1, err2, err3 error
	var statusCode int

	statusCode, err1 = httpProcessGet(url, name, statusList, skipTLSVerify)
	if err1 == nil {
		return statusCode, nil
	}

	statusCode, err2 = httpProcessGet(url, name, statusList, skipTLSVerify)
	if err2 == nil {
		return statusCode, nil
	}

	statusCode, err3 = httpProcessGet(url, name, statusList, skipTLSVerify)
	if err3 == nil {
		return statusCode, nil
	}

	// 去除重复
	var errMap = make(map[string]bool, 3)
	errMap[err1.Error()] = true
	errMap[err2.Error()] = true
	errMap[err3.Error()] = true

	var errStrBuilder strings.Builder
	var n = 0
	for err, _ := range errMap {
		n += 1
		errStrBuilder.WriteString(fmt.Sprintf("检查 %s 错误[%d]: %s; ", name, n, err))
	}

	err := fmt.Errorf("%s", strings.TrimSpace(errStrBuilder.String()))
	return -1, err
}

func httpProcessGet(url string, name string, statusList []string, skipTLSVerify bool) (int, error) {
	// 创建一个自定义的Transport，这样我们可以访问TLS连接状态
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify}, // 忽略服务器证书验证
	}

	// 使用自定义的Transport创建一个HTTP客户端
	client := &http.Client{Transport: tr}

	// 发送请求
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("获取 GET %s 请求错误：%s", name, err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	statusCode := resp.StatusCode

	for _, s := range statusList {
		switch s {
		case "xxx":
			return statusCode, nil
		case "1xx":
			if statusCode >= 100 && statusCode <= 199 {
				return statusCode, nil
			}
		case "2xx":
			if statusCode >= 200 && statusCode <= 299 {
				return statusCode, nil
			}
		case "3xx":
			if statusCode >= 300 && statusCode <= 399 {
				return statusCode, nil
			}
		case "4xx":
			if statusCode >= 400 && statusCode <= 499 {
				return statusCode, nil
			}
		case "5xx":
			if statusCode >= 500 && statusCode <= 599 {
				return statusCode, nil
			}
		default:
			status := fmt.Sprintf("%d", statusCode)
			if status == s {
				return statusCode, nil
			}
		}
	}

	return statusCode, fmt.Errorf("检查 GET %s 状态码错误: %d", name, statusCode)
}
