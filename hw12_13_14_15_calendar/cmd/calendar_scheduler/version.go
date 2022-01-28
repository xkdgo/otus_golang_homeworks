package main

import "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/cmd/calendar_scheduler/service"

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func transferVersionToScheduler() {
	service.Release = release
	service.BuildDate = buildDate
	service.GitHash = gitHash
}
