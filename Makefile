.PHONY: watch start restart

watch:
	@supervisor app.coffee

start:
	@forever start -o /home/joe/rasp-tv.log -e /home/joe/rasp-tv-error.log app.js

restart:
	@forever restart -o /home/joe/rasp-tv.log -e /home/joe/rasp-tv-error.log app.js