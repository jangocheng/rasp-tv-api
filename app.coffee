express = require 'express'
path    = require 'path'
http    = require 'http'
index   = require './routes'
api     = require './routes/api'

app = express()

errorHandler = (err, req, res, next) ->
	console.error err
	res.send 500,
		msg : err.message

pageNotFound = (req, res, next) ->
	res.render '404',
		title : 'Page Not Found'

app.set 'port', 8080
app.set 'views', __dirname + '/views'
app.set 'view engine', 'jade'
app.use express.favicon()
app.use express.bodyParser()
app.use app.router
app.use express.static(path.join(__dirname, 'assets'))
app.use errorHandler
app.use pageNotFound

app.get '/', index.index
app.get '/movies', api.movies
app.post '/movies/play', api.play

server = http.createServer app
io = require('socket.io').listen server
server.listen app.get('port')

console.log 'Listening on port ' + app.get('port')