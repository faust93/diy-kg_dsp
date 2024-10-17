#go build -ldflags="-s -w" main.go oled.go fonts.go
go run main.go oled.go fonts.go images.go rotary.go menu.go cdsp.go config.go
#go build -ldflags="-linkmode=external -extldflags=-static -s -w" main.go oled.go fonts.go images.go rotary.go menu.go cdsp.go config.go
