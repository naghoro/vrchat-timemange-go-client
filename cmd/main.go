package main

import (
	"flag"
	"os"
	"time"
	"vrchat-timemanage/internal/timemanage"
)

func main() {

	duration := flag.Duration("sec", 3600*time.Second, "second of hour")
	flag.Parse()

	manager := timemanage.New(
		timemanage.SetSecondOfHour(*duration),
	)
	os.Exit(manager.ManageStart())
}
