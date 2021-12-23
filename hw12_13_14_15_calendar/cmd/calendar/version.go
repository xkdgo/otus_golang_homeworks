package main

import "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/cmd/calendar/cmd"

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func transferVersionToCalendar() {
	cmd.Release = release
	cmd.BuildDate = buildDate
	cmd.GitHash = gitHash
}
