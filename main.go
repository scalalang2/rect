package main

import "github.com/scalalang2/load-balancing-simulator/reporter"

func main() {
	done := make(chan bool)
	go reporter.ReportAvgStd(done)
	<-done
}