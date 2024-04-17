package aggregator

import (
	"newsaggr/pkg/database/model"
	"newsaggr/pkg/logger"
	"newsaggr/pkg/rss"
	"strings"
	"sync"
	"time"
)

// Aggregator - структура для агрегатора
type Aggregator struct {
	Response      chan []*model.News
	ErrorResponse chan error
	xmlData       *rss.Data
	m             sync.Mutex
}

// New - создание экземпляра агрегатора
func New(xmlData *rss.Data) *Aggregator {
	response := make(chan []*model.News)
	errResponse := make(chan error)
	aggr := &Aggregator{response, errResponse, xmlData, sync.Mutex{}}

	return aggr
}

func waiter(dur int, done chan struct{}) {
	time.Sleep(time.Minute * time.Duration(dur))
	done <- struct{}{}
}

func (a *Aggregator) handleWriteToBase(news []*model.News) {
	for _, data := range news {
		if err := data.Create(); err != nil {
			a.ErrorResponse <- err
			continue
		}
	}
}

func (a *Aggregator) handleRounder(url string) {
	data, err := rss.Round(url)
	if err != nil {
		a.ErrorResponse <- err
		return
	}

	var news []*model.News

	for _, item := range data.Channel.Item {
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		data := &model.News{
			Title:   item.Title,
			Content: item.Description,
			Link:    item.Link,
		}
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			data.PubTime = t.Unix()
		}
		news = append(news, data)
	}

	a.Response <- news
}

func (a *Aggregator) handler() {
	done := make(chan struct{})
	for {
		for _, url := range a.xmlData.URLS {
			go a.handleRounder(url)
		}
		go waiter(int(a.xmlData.RequestPeriod), done)
		<-done
	}
}

func (a *Aggregator) responseWorker() {
	for {
		select {
		case news := <-a.Response:
			go a.handleWriteToBase(news)
			logger.Info("Получили новость")
		case err := <-a.ErrorResponse:
			logger.Error("Ошибка: ", err)
		default:
			logger.Info("Ожидание")
		}
	}
}

func (a *Aggregator) Start() {
	go a.responseWorker()
	go a.handler()
}
