express = require 'express'
path    = require 'path'

app = express()

errorPage = (err, req, res, next) ->
    console.error 'Error: ' + err.message
    res.send 500,
        error : err.message

pageNotFound = (req, res, next) ->
    res.render '404',
        title : 'Page Not Found'

app.configure () ->
    app.set 'port', 8080
    app.set 'views', __dirname + '/views'
    app.set 'view engine', 'jade'
    app.use express.favicon()
    app.use express.bodyParser()
    app.use app.router
    app.use express.static(path.join(__dirname, 'assets'))
    app.use errorPage
    app.use pageNotFound

index = require './routes'

app.get '/', index.index

app.listen app.get('port')
console.log 'Listening on port ' + app.get('port')