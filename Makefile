JSDIR=static/js
EXCLUDES= \
	--exclude=".git*" \
	--exclude="LICENSE" \
	--exclude="README.md" \
	--exclude="rasp-tv-api" \
	--exclude="raspTv.db" \
	--exclude="logs.txt" \
	--exclude="Makefile"

GOSOURCE=$(shell find . -type f -name "*.go")

.PHONY: clean watch

all: rasp-tv

clean:
	rm -fr rasp-tv-api dist

rasp-tv: $(GOSOURCE)
	go build

deploy:
	rsync -avz --delete $(EXCLUDES) ./ ./dist
	rsync -avz --delete ./dist/ joe@192.168.11.16:/home/joe/workspace/go/src/github.com/simonjm/rasp-tv-api
	$(MAKE) clean
