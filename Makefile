.PHONY: all
all: docs/spec.txt docs/tokens4.txt docs/tokens-all.json docs/tokens-uniq.txt docs/word2stem.txt docs/word2stem.json docs/count.txt docs/uniq.txt docs/dic.ja.json web

.PHONY: data/dic.ja.tsv
data/dic.ja.tsv:
	wget -O data/dic.ja.tsv 'https://docs.google.com/spreadsheets/d/e/2PACX-1vSLRyGpO5qAUt2YGejK3tkELmnrGKHX0iALEFIgdN0vKCOZU0j9lseDLf3s8UA8waZnL3uAWsDk1Xp7/pub?gid=406497718&single=true&output=tsv'

bin/s2t: spec2text/* spec2text/*/*
	go build -o $@ ./spec2text/cmd

docs/spec.txt: spec.html bin/s2t
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

docs/word2stem.txt: docs/tokens-uniq.txt gospec
	./gospec normalize < $< > $@ 2> docs/word2stem.log
	cat docs/word2stem.log | sort | uniq > docs/word2stem.uniq.log

docs/word2stem.json: docs/tokens-uniq.txt gospec
	./gospec normalizejson < $< > $@ 2>/dev/null

docs/word2stem.js: docs/word2stem.json
	echo 'var word2stem = ' > $@
	cat $< >> $@

docs/count.txt: docs/tokens4.txt gospec
	./gospec count < $< > $@ 2>/dev/null

docs/uniq.txt: docs/tokens4.txt gospec
	./gospec uniq < $< > $@ 2>/dev/null

docs/dic.ja.json: data/dic.ja.tsv bin/tsv2json
	bin/tsv2json $< > $@

docs/dic.ja.js: docs/dic.ja.json
	echo 'var dic = ' > $@
	cat $< >> $@

.PHONEY: web
web: docs/spec.html docs/lib/godoc/style.css docs/main.js docs/dic.ja.js docs/word2stem.js docs/lib/godoc/jquery.js docs/lib/godoc/playground.js docs/lib/godoc/godocs.js

spec_noscript.html: spec_orig.html
	perl -p -e 'BEGIN{undef $$/;}  s|<script>[^<]*</script>||smg' $< > $@

docs/spec.html: spec.html
	mkdir -p docs
	cp $< $@

docs/lib/godoc/style.css: lib/godoc/style.css
	cp $< $@

docs/dictionary.css: dictionary.css
	cp $< $@

docs/main.js: main.js
	cp $< $@

bin/tsv2json: tsv2json/*/*
	go mod vendor
	go build -o $@ ./tsv2json/cmd

docs/lib/godoc/jquery.js:
	touch $@

docs/lib/godoc/playground.js:
	touch $@

docs/lib/godoc/godocs.js:
	touch $@
