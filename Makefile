GOBIN = ./build/bin
GO ?= latest
GOBUILD = env GO111MODULE=on go build

opcodeAvg:
	$(GOBUILD) -o $(GOBIN)/opcodeAvg ./cmd/opcodeAvg/
	@echo "Done building."
	@echo "Run "$(GOBIN)/opcodeAvg" to start"

sacc:
	$(GOBUILD) -o $(GOBIN)/sacc ./cmd/sacc/
	@echo "Done building."
	@echo "Run "$(GOBIN)/sacc" to start"

garet:
	$(GOBUILD) -o $(GOBIN)/sacc ./cmd/garet/
	@echo "Done building."
	@echo "Run "$(GOBIN)/garet" to start"

garet:
	$(GOBUILD) -o $(GOBIN)/balanceMeter ./cmd/balanceMeter/
	@echo "Done building."
	@echo "Run "$(GOBIN)/balanceMeter" to start"