VERSION = 0.1.$(shell date +%Y%m%d.%H%M)
FLAGS := "-s -w -X main.version=${VERSION}"

all: build readme

build:
	CGO_ENABLED=0 go build -ldflags=${FLAGS} -o http-cert-provisioner .

readme:
	cp HEAD.md README.md
	echo -e '\n```\n# http-cert-provisioner -h' >> README.md
	./http-cert-provisioner -h 2>> README.md
	echo -e '```' >> README.md
