.PHONY: all

all: docs/spec.txt docs/tokens.txt docs/count.txt docs/uniq.txt

gospec: *.go
	go build

docs/spec.txt: spec2text/*.go spec.html
	./gospec text > docs/spec.txt

docs/tokens.txt: gospec docs/spec.txt
	./gospec dump > docs/tokens.txt

docs/count.txt: gospec docs/spec.txt
	./gospec count > docs/count.txt

docs/uniq.txt: gospec docs/spec.txt
	./gospec uniq > docs/uniq.txt
