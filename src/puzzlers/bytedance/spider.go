package main

import (
	"fmt"
	"net/http"
	"time"
)

type crawlResult struct {
	err        error
	httpStatus string
	channel    string
}

func crawlChannel(channel string, ch chan crawlResult) {
	time.Sleep(time.Second * 5)
	resp, err := http.Get(channel)
	if err != nil {
		ch <- crawlResult{err: fmt.Errorf("crawl failed, %s", channel, err), channel: channel}
		return
	}
	ch <- crawlResult{httpStatus: resp.Status, channel: channel}
}

func refresh(gCh, cCh chan bool) {
	const noticeStr = `-\|/`
	i := 0
	for {
		fmt.Printf("\r%c", noticeStr[i])
		i += 1
		i %= len(noticeStr)
		select {
		case <-gCh:
			cCh <- true
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func main() {
	ttChannels := [...]string{
		"http://www.toutiao.com/ch/news_hot/",
		"http://www.toutiao.com/ch/essay_joke/",
		"http://www.toutiao.com/ch/news_society/",
		"http://test.test.com/",
	}
	workerCh := make(chan crawlResult, len(ttChannels))
	for _, channel := range ttChannels {
		go crawlChannel(channel, workerCh)
	}

	gCh := make(chan bool, 1)
	cCh := make(chan bool, 1)
	go refresh(gCh, cCh)

	crawlResults := make(map[string]crawlResult)
	for i := 0; i < len(ttChannels); i++ {
		result := <-workerCh
		crawlResults[result.channel] = result
	}

	// finished, stop refreshing
	gCh <- true
	<-cCh

	for k, v := range crawlResults {
		fmt.Printf("\r%s, ", k)
		if v.err != nil {
			fmt.Printf("%s", v.err)
		}
		fmt.Printf("httpStatus is %s", v.httpStatus)
		fmt.Printf("\n")
	}
}
