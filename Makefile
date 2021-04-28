.PHONY: all
all: spec2text docs/spec.txt docs/tokens.txt docs/count.txt docs/uniq.txt

s2t: spec.html spec2text/*/*
	go build -o s2t ./spec2text/cmd

docs/spec.txt: spec.html
	./s2t spec.html > docs/spec.txt

gospec: *.go
	go build

docs/tokens.txt: gospec docs/spec.txt
	./gospec dump > docs/tokens.txt

docs/count.txt: gospec docs/spec.txt
	./gospec count > docs/count.txt

docs/uniq.txt: gospec docs/spec.txt
	./gospec uniq > docs/uniq.txt
