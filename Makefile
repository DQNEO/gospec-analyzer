# force GNU sed
ifeq ($(shell uname),Linux)
	# for Linux
	SED := sed
else
	# for MacOS assuming GNU sed is installed as "gsed"
	SED := gsed
endif

.PHONY: all
all: docs/spec.txt docs/tokens4.txt docs/tokens-all.json docs/tokens-uniq.txt docs/word2lemma.txt docs/word2lemma.json docs/count_by_lemma_and_tag.txt docs/count_by_lemma.txt web

.PHONY: spec_orig.html
spec_orig.html:
	wget -O spec_orig.html 'https://tip.golang.org/ref/spec'

bin/s2t: spec2text/* spec2text/*/*
	go build -o $@ ./spec2text/cmd

docs/spec.txt: spec_orig.html bin/s2t
	@echo --- extracting text from spec_orig.html
	bin/s2t $< > $@

bin/tokenizer: tokenizer/* tokenizer/*/*
	go build -o $@ ./tokenizer/cmd

docs/tokens-all.txt: docs/spec.txt bin/tokenizer
	bin/tokenizer $< > $@

docs/tokens-all.json: docs/spec.txt bin/tokenizer
	bin/tokenizer --json $< > $@

gospec: *.go
	go build -o gospec .

docs/tokens1.txt: docs/tokens-all.txt gospec
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

docs/tokens-uniq.txt: docs/tokens4.txt
	cat $< | sort | uniq | tr '[:upper:]' '[:lower:]' > $@

docs/word2lemma.txt: docs/tokens-uniq.txt gospec
	./gospec lemmatize < $< > $@ 2> docs/word2lemma.log
	cat docs/word2lemma.log | sort | uniq > docs/word2lemma.uniq.log

docs/word2lemma.json: docs/tokens-uniq.txt gospec
	./gospec lemmatizejson < $< > $@ 2>/dev/null

docs/count_by_lemma_and_tag.txt: docs/tokens4.txt gospec
	./gospec count_by_lemma_and_tag < $< > $@ 2>/dev/null

docs/count_by_lemma.txt: docs/tokens4.txt gospec
	./gospec count_by_lemma < $< > $@ 2>/dev/null

bin/tsv2json: tsv2json/*/*
	go mod vendor
	go build -o $@ ./tsv2json/cmd

.PHONY: web
web: docs/word2lemma.json
