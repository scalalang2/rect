GOBIN = ./build/bin
GO ?= latest
GOBUILD = env GO111MODULE=on go build

opcodeAvg:
	$(GOBUILD) -o $(GOBIN)/opcodeAvg ./cmd/opcodeAvg/
	@echo "Done building."
	@echo "Run "$(GOBIN)/opcodeAvg" to start"