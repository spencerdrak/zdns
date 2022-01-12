module github.com/zmap/zdns

go 1.17

require (
	github.com/hashicorp/go-version v1.2.0
	github.com/liip/sheriff v0.0.0-20190308094614-91aa83a45a3d
	github.com/miekg/dns v1.1.27
	github.com/sirupsen/logrus v1.4.2
	github.com/zmap/go-iptree v0.0.0-20170831022036-1948b1097e25
)

require (
	github.com/asergeyev/nradix v0.0.0-20170505151046-3872ab85bb56 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)

replace github.com/miekg/dns => github.com/zmap/dns v1.1.35-zdns-2
