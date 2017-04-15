########################################################################################

# This Makefile generated by GoMakeGen 0.4.0 using next command:
# gomakegen .

########################################################################################

.PHONY = fmt all clean deps

########################################################################################

all: deadline

deadline:
	go build deadline.go

deps:
	git config --global http.https://pkg.re.followRedirects true
	go get -v pkg.re/essentialkaos/ek.v8

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

clean:
	rm -f deadline

########################################################################################
