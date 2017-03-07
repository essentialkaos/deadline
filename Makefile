########################################################################################

.PHONY = fmt all clean deps

########################################################################################

all: deadline

deadline:
	go build deadline.go

deps:
	go get -v pkg.re/essentialkaos/ek.v7

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

clean:
	rm -f deadline

########################################################################################

