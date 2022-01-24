package main

import "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/cmd/calendar_scheduler/service"

func main() {
	transferVersionToScheduler()
	service.Execute()
}
