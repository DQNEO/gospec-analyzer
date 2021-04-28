.PHONY: all
all: docs/spec.txt docs/tokens.tsv # docs/count.txt docs/uniq.txt

s2t: spec.html spec2text/*/*
	go build -o s2t ./spec2text/cmd

docs/spec.txt: spec.html
	./s2t spec.html > docs/spec.txt

prs: prose/*/*
	go build -o prs ./prose/cmd

docs/tokens.tsv: prs docs/spec.txt
	./prs docs/spec.txt > docs/tokens.tsv

docs/tokens.json: prs docs/spec.txt
	./prs --json docs/spec.txt > docs/tokens.json

gospec: *.go
	go build -o gospec .

docs/count.txt: gospec docs/spec.txt
	./gospec count < docs/spec.txt> docs/count.txt

docs/uniq.txt: gospec docs/spec.txt
	./gospec uniq < docs/spec.txt > docs/uniq.txt
