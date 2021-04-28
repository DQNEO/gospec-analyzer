.PHONY: all
all: s2t docs/spec.txt docs/tokens.txt docs/count.txt docs/uniq.txt

s2t: spec.html spec2text/*/*
	go build -o s2t ./spec2text/cmd

docs/spec.txt: spec.html
	./s2t spec.html > docs/spec.txt

gospec: *.go
	go build -o gospec .

docs/tokens.txt: gospec docs/spec.txt
	./gospec dump < docs/spec.txt > docs/tokens.txt

docs/count.txt: gospec docs/spec.txt
	./gospec count < docs/spec.txt> docs/count.txt

docs/uniq.txt: gospec docs/spec.txt
	./gospec uniq < docs/spec.txt > docs/uniq.txt
