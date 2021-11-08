.PHONY: buildv2012
BIN_FILE=GBT19056Gen.exe

buildv2012:
		@go build -o target/GBT19056-2012Gen.exe main.go

buildv2020:
		@go build -o target/GBT19056-2020Gen.exe main2020.go