JSDIR=assets/js
EXCLUDES=--exclude=".git*" --exclude="assets/js/*.coffee" --exclude="LICENSE" --exclude="README.md" --exclude="rasp-tv" --exclude="raspTv.db" --exclude="logs.txt"
JSLIBS=$(JSDIR)/libs/angular.js $(JSDIR)/libs/angular-animate.js $(JSDIR)/libs/angular-route.js $(JSDIR)/libs/angular-resource.js $(JSDIR)/libs/loading-bar.js $(JSDIR)/libs/angular-strap.js $(JSDIR)/libs/angular-strap.tpl.js $(JSDIR)/libs/filter-regex.js

.PHONY: clean watch

all: $(JSDIR)/app.js $(JSDIR)/services.js rasp-tv

clean:
	rm -fr rasp-tv assets/js/*.js dist

$(JSDIR)/%.js: $(JSDIR)/%.coffee
	coffee -c $<

watch:
	coffee -wc assets/js/*.coffee

rasp-tv: app.go api/helpers.go api/index.go api/movies.go api/scan.go api/shows.go api/player.go data/movies.go data/shows.go data/sessions.go
	go build

deploy: $(JSDIR)/rasptv.min.js
	rsync -avz --delete $(EXCLUDES) ./ ./dist
	rsync -avz --delete ./dist/ joe@192.168.11.2:/home/joe/go/src/simongeeks.com/joe/rasp-tv
	# @ssh joe@rpi "cd /home/joe/go/src/simongeeks.com/joe/rasp-tv && make rasp-tv"
	$(MAKE) clean

$(JSDIR)/rasptv.min.js: $(JSLIBS) $(JSDIR)/services.js $(JSDIR)/app.js
	cat $^ | uglifyjs -c -m --screw-ie8 -o $@ -
