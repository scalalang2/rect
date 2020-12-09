package main

import "load-balancing-simulator/reporter"

func main() {
	done := make(chan bool)
	go reporter.ReportAvgStd(done)
	<-done
}