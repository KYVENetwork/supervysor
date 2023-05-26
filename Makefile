#!/usr/bin/make -f

supervysor:
	@echo "🤖 Building supervysor..."
	@go build -mod=readonly -o ./build/supervysor ./cmd/supervysor/main.go
	@echo "🤖️Copy into ~/go/bin/supervysor..."
	@cp build/supervysor ~/go/bin/supervysor
	@echo "✅ Completed build!"
