version := $(shell /bin/date "+%Y%m%d")
module_name := cymzjs
module_type := api

release_name := $(module_name)-$(module_type)

linux:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o $(release_name) $(module_name).go
	command -v upx &> /dev/null && upx $(release_name)

	scp $(release_name) root@xg1:/data/site/grj/$(release_name)
	rm -rf $(release_name)

