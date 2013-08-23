fs   = require 'fs'
path = require 'path'
omx  = require './lib/omxcontrol'

moviesDir = '/media/passport/Movies'
# moviesDir = '/Users/Joe/Movies/Test'
showsDir  = '/media/passport/TV Shows'

exports.movies = (req, res, next) ->
    fs.readdir moviesDir, (err, files) ->
        next(err) if err
        movies = files.filter (file) ->
            file[0] isnt '.'
        res.json
            'movies' : movies.sort()

exports.play = (req, res, next) ->
    moviePath = path.join moviesDir, req.body.movie
    fs.readdir moviePath, (err, files) ->
        next(err) if err
        movieFiles = files.filter (file) ->
            path.extname(file) in ['.mp4', '.avi', '.mov', '.mkv']
        if movieFiles.length > 0
            omx.quit()
            omx.start path.join(moviePath, movieFiles[0])
            # console.log 'Movie playing from ' + path.join moviePath, movieFiles[0]
        else
            next(new Error 'Movie file not found')
        res.send 200