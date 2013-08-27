.PHONY: watch start restart

watch:
	@supervisor app.coffee

start:
	@NODE_ENV=production forever start -o /home/joe/rasp-tv.log -e /home/joe/rasp-tv-error.log app.js

restart:
	@NODE_ENV=production forever restart -o /home/joe/rasp-tv.log -e /home/joe/rasp-tv-error.log app.js