.PHONY: all
all: docs/spec.txt docs/tokens.txt docs/tokens.json  docs/count.txt docs/uniq.txt

bin/s2t: spec.html spec2text/*/*
	go build -o bin/s2t ./spec2text/cmd

docs/spec.txt: bin/s2t spec.html
	bin/s2t spec.html > docs/spec.txt

bin/prose: prose/*/*
	go build -o bin/prose ./prose/cmd

docs/tokens.txt: bin/prose docs/spec.txt
	bin/prose docs/spec.txt > docs/tokens.txt

docs/tokens.json: bin/prose docs/spec.txt
	bin/prose --json docs/spec.txt > docs/tokens.json

gospec: *.go
	go build -o gospec .

docs/count.txt: gospec docs/spec.txt
	./gospec count < docs/tokens.txt > docs/count.txt

docs/uniq.txt: gospec docs/spec.txt
	./gospec uniq < docs/tokens.txt > docs/uniq.txt
