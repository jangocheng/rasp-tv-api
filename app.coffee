express  = require 'express'
path     = require 'path'
http     = require 'http'
index    = require './routes'
movies   = require './routes/movies'
shows    = require './routes/shows'
player   = require './routes/player'
shutdown = require './routes/shutdown'
youtube  = require './routes/youtube'

process.env.NODE_ENV = if process.argv[2]? then 'production' else null

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

app.configure 'production', () ->
    app.use express.compress()
    app.locals.production = true

app.use express.favicon()
app.use express.json()
app.use express.urlencoded()
app.use app.router
app.use express.static(path.join(__dirname, 'assets'), {maxAge : 31557600000})
app.use errorHandler
app.use pageNotFound

app.get '/', index.index
app.get '/movies', movies.index
app.post '/movies/play', movies.play
app.get '/shows', shows.index
app.post '/shows/seasons', shows.seasons
app.post '/shows/seasons/episodes', shows.episodes
app.post '/shows/play', shows.play
app.post '/shutdown', shutdown
app.post '/shows/random', shows.random
app.get '/youtube/videos', youtube.get
app.post '/youtube/play', youtube.play

server = http.createServer app
io = require('socket.io').listen server
server.listen app.get('port')
io.sockets.on 'connection', player

console.log 'Listening on port ' + app.get('port')