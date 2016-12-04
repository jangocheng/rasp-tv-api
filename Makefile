JSDIR=static/js
EXCLUDES= \
	--exclude=".git*" \
	--exclude="static/js/*.coffee" \
	--exclude="LICENSE" \
	--exclude="README.md" \
	--exclude="rasp-tv" \
	--exclude="raspTv.db" \
	--exclude="logs.txt"

JSLIBS= \
	$(JSDIR)/libs/angular.js \
	$(JSDIR)/libs/angular-animate.js \
	$(JSDIR)/libs/angular-route.js \
	$(JSDIR)/libs/angular-resource.js \
	$(JSDIR)/libs/loading-bar.js \
	$(JSDIR)/libs/angular-strap.js \
	$(JSDIR)/libs/angular-strap.tpl.js \
	$(JSDIR)/libs/filter-regex.js

GOSOURCE=$(shell find . -type f -name "*.go")

.PHONY: clean watch

all: $(JSDIR)/app.js $(JSDIR)/services.js rasp-tv

clean:
	rm -fr rasp-tv static/js/*.js dist

$(JSDIR)/%.js: $(JSDIR)/%.coffee
	coffee -c $<

watch:
	coffee -wc static/js/*.coffee

rasp-tv: $(GOSOURCE)
	go build

deploy: $(JSDIR)/rasptv.min.js
	rsync -avz --delete $(EXCLUDES) ./ ./dist
	# rsync -avz --delete ./dist/ joe@192.168.11.16:/home/joe/workspace/go/src/simongeeks.com/joe/rasp-tv
	# $(MAKE) clean

$(JSDIR)/rasptv.min.js: $(JSLIBS) $(JSDIR)/services.js $(JSDIR)/app.js
	cat $^ | uglifyjs -c -m --screw-ie8 -o $@ -
