all: bootstrap precheck test

FILES := $$(find . -name '*.go' | grep -vE 'vendor')
SOURCE := .

golint:
	go get golang.org/x/lint/golint

godep:
	go get github.com/tools/godep

staticcheck:
	go get honnef.co/go/tools/cmd/staticcheck

revive:
	go get github.com/mgechev/revive

errcheck:
	go get github.com/kisielk/errcheck

check:
	go get gitlab.com/opennota/check/cmd/aligncheck
	go get gitlab.com/opennota/check/cmd/structcheck
	go get gitlab.com/opennota/check/cmd/varcheck

bootstrap: golint errcheck check staticcheck

precheck:
	# format all go files
	gofmt -l -s -w $(SOURCE)/

	# golint all go files
	@for file in $(FILES); do golint $$file; done;

	# aligncheck - net align check
	# For the visualisation of struct packing see http://golang-sizeof.tips/
	aligncheck $(SOURCE)/...

	# Find unused struct fields.
	structcheck $(SOURCE)/...

	# Find unused global variables and constants.
	varcheck $(SOURCE)/...

	# error check
	errcheck ./$(SOURCE)/...

	# Vet examines Go source code and reports suspicious constructs
	go vet -vettool=$(which shadow) ./... 2>&1

	staticcheck -go 1.13 -tests ./...

test: precheck
	# run test
	go test $(SOURCE)/...

