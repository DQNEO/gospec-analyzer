.PHONY: all
all: docs/spec.txt docs/tokens2.txt docs/tokens0.json docs/normalized.txt  docs/count.txt docs/uniq.txt

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
	./gospec filter1 < docs/tokens0.txt > docs/tokens1.txt 2> docs/tokens1.log
	cat docs/tokens1.log | sort | uniq > docs/tokens1.uniq.log

docs/tokens2.txt: gospec docs/tokens1.txt
	./gospec filter2 < docs/tokens1.txt > docs/tokens2.txt 2> docs/tokens2.log
	cat docs/tokens2.log | sort | uniq > docs/tokens2.uniq.log

docs/tokens3.txt: gospec docs/tokens2.txt
	./gospec filter3 < docs/tokens2.txt > docs/tokens3.txt 2> docs/tokens3.log
	cat docs/tokens3.log | sort | uniq > docs/tokens3.uniq.log

docs/tokens4.txt: gospec docs/tokens3.txt
	./gospec filter4 < docs/tokens3.txt > docs/tokens4.txt 2> docs/tokens4.log
	cat docs/tokens4.log | sort | uniq > docs/tokens4.uniq.log

docs/normalized.txt: gospec docs/tokens4.txt
	./gospec normalize < docs/tokens4.txt > docs/normalized.txt 2> docs/normalized.log
	cat docs/normalized.log | sort | uniq > docs/normalized.uniq.log

docs/count.txt: gospec docs/tokens4.txt
	./gospec count < docs/tokens4.txt > docs/count.txt 2>/dev/null

docs/uniq.txt: gospec docs/tokens4.txt
	./gospec uniq < docs/tokens4.txt> docs/uniq.txt 2>/dev/null
