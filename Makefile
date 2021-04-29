.PHONY: all
all: docs/spec.txt docs/tokens2.txt docs/tokens0.json  docs/count.txt docs/uniq.txt

bin/s2t: spec.html spec2text/*/*
	go build -o bin/s2t ./spec2text/cmd

docs/spec.txt: bin/s2t spec.html
	bin/s2t spec.html > docs/spec.txt

bin/tokenizer: tokenizer
	go build -o bin/tokenizer ./tokenizer/cmd

docs/tokens0.txt: bin/tokenizer docs/spec.txt
	bin/tokenizer docs/spec.txt > docs/tokens0.txt

docs/tokens0.json: bin/tokenizer docs/spec.txt
	bin/tokenizer --json docs/spec.txt > docs/tokens0.json

gospec: *.go
	go build -o gospec .

docs/tokens1.txt: gospec docs/tokens0.txt
	./gospec filter1 < docs/tokens0.txt > docs/tokens1.txt

docs/tokens2.txt: gospec docs/tokens1.txt
	./gospec filter2 < docs/tokens1.txt > docs/tokens2.txt

docs/tokens3.txt: gospec docs/tokens2.txt
	./gospec filter2 < docs/tokens2.txt > docs/tokens3.txt

docs/count.txt: gospec docs/tokens3.txt
	./gospec count < docs/tokens3.txt > docs/count.txt

docs/uniq.txt: gospec docs/tokens3.txt
	./gospec uniq < docs/tokens3.txt> docs/uniq.txt
