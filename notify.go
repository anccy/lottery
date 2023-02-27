package main

import (
	"context"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/lark"
	"log"
	"strconv"
	"time"
)

func NotifyLark(webHookURL, subject, message string) {
	larkWebhookSvc := lark.NewWebhookService(webHookURL)

	notifier := notify.New()
	notifier.UseServices(larkWebhookSvc)

	if err := notifier.Send(context.Background(), subject, message); err != nil {
		log.Fatalf("notifier.Send() failed: %s", err.Error())
	}

	log.Println("notification sent")
}

func Cycle() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case now := <-ticker.C:
			// check if today is Tuesday, Thursday, or Sunday and the current time is 18:00
			for day := range Config.DaysOfWeek {
				if int(now.Weekday()) == day {
					if now.Hour() == 18 && now.Minute() == 0 && now.Second() == 0 {
						data, err := GetLotteryData(100)
						if err != nil {
							log.Fatal(err)
						}
						ret := PredictByMode(data, Config.PredictNum)
						retStr := strconv.Itoa(ret[0])
						for i := 1; i < len(ret); i++ {
							retStr = retStr + " " + strconv.Itoa(ret[i])
						}
						NotifyLark(Config.WebHookURL, "今日财富密码", retStr)
					}
				}
			}
		}
	}
}
