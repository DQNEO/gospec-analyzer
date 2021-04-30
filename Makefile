.PHONY: all
all: docs/spec.txt docs/tokens2.txt docs/tokens0.json docs/normalized.txt  docs/count.txt docs/uniq.txt

bin/s2t: spec2text/*/* spec.html
	go build -o $@ ./spec2text/cmd

docs/spec.txt: spec.html bin/s2t
	bin/s2t $< > $@

bin/tokenizer: tokenizer/*/*
	go build -o $@ ./tokenizer/cmd

docs/tokens0.txt: docs/spec.txt bin/tokenizer
	bin/tokenizer $< > $@

docs/tokens0.json: docs/spec.txt bin/tokenizer
	bin/tokenizer --json $< > $@

gospec: *.go
	go build -o gospec .

docs/tokens1.txt: docs/tokens0.txt gospec
	./gospec filter1 < $< > $@ 2> docs/tokens1.log
	cat docs/tokens1.log | sort | uniq > docs/tokens1.uniq.log

docs/tokens2.txt: docs/tokens1.txt gospec
	./gospec filter2 < $< > $@ 2> docs/tokens2.log
	cat docs/tokens2.log | sort | uniq > docs/tokens2.uniq.log

docs/tokens3.txt: docs/tokens2.txt gospec
	./gospec filter3 < $< > $@ 2> docs/tokens3.log
	cat docs/tokens3.log | sort | uniq > docs/tokens3.uniq.log

docs/tokens4.txt: docs/tokens3.txt gospec
	./gospec filter4 < $< > $@ 2> docs/tokens4.log
	cat docs/tokens4.log | sort | uniq > docs/tokens4.uniq.log

docs/normalized.txt: docs/tokens4.txt gospec
	./gospec normalize < $< > $@ 2> docs/normalized.log
	cat docs/normalized.log | sort | uniq > docs/normalized.uniq.log

docs/count.txt: docs/tokens4.txt gospec
	./gospec count < $< > $@ 2>/dev/null

docs/uniq.txt: docs/tokens4.txt gospec
	./gospec uniq < $< > $@ 2>/dev/null
