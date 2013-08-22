fs   = require 'fs'
path = require 'path'
omx  = require './lib/omxcontrol'

# file = '/Users/Joe/Movies/Test'
file = '/media/passport/'

exports.movies = (req, res, next) ->
    moviesDir = path.join file, 'Movies'
    fs.readdir moviesDir, (err, files) ->
        next(err) if err
        movies = files.filter (file) ->
            file[0] isnt '.'
        res.json
            'movies' : movies.sort()

exports.play = (req, res, next) ->
    moviePath = path.join file, 'Movies', req.body.movie
    fs.readdir moviePath, (err, files) ->
        next(err) if err
        movieFiles = files.filter (file) ->
            path.extname(file) in ['.mp4', '.avi', '.mov']
        if movieFiles.length > 0
            omx.quit()
            omx.start path.join(moviePath, movieFiles[0])
        else
            next(new Error 'Movie file not found')
        res.send 200