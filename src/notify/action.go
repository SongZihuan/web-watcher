package notify

import (
	"fmt"
	"github.com/SongZihuan/web-watcher/src/config"
	"github.com/SongZihuan/web-watcher/src/smtpserver"
	"github.com/SongZihuan/web-watcher/src/wxrobot"
	"strings"
	"sync"
	"time"
)

type urlRecord struct {
	Name string
	URL  string
	Err  string
}

var startTime time.Time
var records sync.Map

func InitNotify() error {
	if !config.IsReady() {
		panic("config is not ready")
	}

	startTime = time.Now().In(config.TimeZone())

	err := smtpserver.InitSmtp()
	if err != nil {
		return err
	}

	return nil
}

func NewRecord(name string, url string, err string) {
	if name == "" {
		name = url
	}

	records.Store(name, &urlRecord{
		Name: name,
		URL:  url,
		Err:  err,
	})
}

func SendNotify() {
	var res strings.Builder
	var count uint64 = 0

	res.WriteString(fmt.Sprintf("日期：%s %s\n", startTime.Format("2006-01-02 15:04:05"), startTime.Location().String()))

	records.Range(func(key, value any) bool {
		record, ok := value.(*urlRecord)
		if !ok {
			return true
		}

		count += 1
		res.WriteString(fmt.Sprintf("- %s 异常: %s\n", record.Name, record.Err))

		return true
	})

	if count <= 0 {
		// 无任何记录
		return
	}

	res.WriteString(fmt.Sprintf("共计：%d 条。\n", count))
	res.WriteString("完毕\n")
	msg := res.String()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		wxrobot.SendNotify(msg)
	}()

	go func() {
		defer wg.Done()
		smtpserver.SendNotify(msg)
	}()

	wg.Wait()
}
