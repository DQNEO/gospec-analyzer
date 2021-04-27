.PHONY: all

all: out/spec.txt out/tokens.txt out/count.txt out/uniq.txt

gospec: *.go
	go build

out/spec.txt: spec2text/*.go spec.html
	./gospec text > out/spec.txt

out/tokens.txt: gospec out/spec.txt
	./gospec dump > out/tokens.txt

out/count.txt: gospec out/spec.txt
	./gospec count > out/count.txt

out/uniq.txt: gospec out/spec.txt
	./gospec uniq > out/uniq.txt
