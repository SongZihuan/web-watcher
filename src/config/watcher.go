package config

import (
	"crypto/tls"
	"fmt"
	"github.com/SongZihuan/web-watcher/src/utils"
	"strconv"
)

type URLConfig struct {
	Name          string           `yaml:"name"`
	URL           string           `yaml:"url"`
	SkipTLSVerify utils.StringBool `yaml:"skip-tls-verify"`
	Status        []string         `yaml:"status"`
	ClientCert    string           `yaml:"client-cert"`
	ClientKey     string           `yaml:"client-key"`

	ClientKeyPair *tls.Certificate `yaml:"-"`
}

type WatcherConfig struct {
	URLs []*URLConfig `yaml:"urls"`
}

func (w *WatcherConfig) setDefault() {
	for _, url := range w.URLs {
		if url.Name == "" {
			host, err := utils.GetURLHost(url.URL)
			if err == nil {
				url.Name = host
			} else {
				url.Name = url.URL
			}
		}

		url.SkipTLSVerify.SetDefaultDisable()

		if len(url.Status) == 0 {
			url.Status = []string{"2xx"}
		}
	}
	return
}

func (w *WatcherConfig) check() (err ConfigError) {
	if len(w.URLs) == 0 {
		return NewConfigError("not any urls")
	}

	for _, url := range w.URLs {
		if !utils.IsValidHTTPURL(url.URL) {
			return NewConfigError(fmt.Sprintf("'%s' is not a valid http/https url", url))
		}

	StatusCycle:
		for _, s := range url.Status {
			switch s {
			case "xxx":
				fallthrough
			case "1xx":
				fallthrough
			case "2xx":
				fallthrough
			case "3xx":
				fallthrough
			case "4xx":
				fallthrough
			case "5xx":
				continue StatusCycle
			default:
				sNum, err := strconv.ParseUint(s, 10, 64)
				if err != nil || sNum < 100 || sNum > 599 {
					return NewConfigError(fmt.Sprintf("'%s' is not a valid status code", s))
				}
				continue StatusCycle
			}
		}

		if url.ClientCert != "" && url.ClientKey != "" {
			cert, err := tls.X509KeyPair([]byte(url.ClientCert), []byte(url.ClientKey))
			if err != nil {
				return NewConfigError(fmt.Sprintf("the cert/key of '%s' is not valid", url.URL))
			}
			url.ClientKeyPair = &cert
			fmt.Printf("%s load client tls cert\n", url.URL)
		} else {
			url.ClientKeyPair = nil
		}
	}

	return nil
}
