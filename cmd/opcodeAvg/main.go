package main

import "github.com/scalalang2/load-balancing-simulator/reporter"

func main() {
	doneReportAvgStd := make(chan bool)
	go reporter.ReportAvgStd(doneReportAvgStd)
	<-doneReportAvgStd
}
