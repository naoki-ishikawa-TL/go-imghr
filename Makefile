build:
	go build -ldflags "-X main.Version=`git rev-parse HEAD`" -tags=release imghr.go

run:
	go run imghr.go
	
install-deps:
	go get -u github.com/f110/go-ihr
	go get -u github.com/f110/go-ihr-slack
	go get -u golang.org/x/net/websocket
	go get -u github.com/gographics/imagick/imagick
	go get -u github.com/moovweb/gokogiri
	go get -u github.com/jteeuwen/go-bindata/...
	
asset-pack:
	go-bindata -pkg=assets -o=./assets/assets.go data

clean:
	rm -f imghr