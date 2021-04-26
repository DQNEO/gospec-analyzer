.PHONY: all

all: out/spec.txt out/tokens.txt out/count.txt

gospec: *.go
	go build

out/spec.txt: spec2text/*.go spec.html
	./gospec text > out/spec.txt

out/tokens.txt: gospec spec.html
	./gospec dump > out/tokens.txt

out/count.txt: gospec spec.html
	./gospec count > out/count.txt
