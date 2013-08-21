fs   = require 'fs'
path = require 'path'
omx  = require 'omxcontrol'

file = '/Users/Joe/Movies/Test/'

exports.movies = (req, res, next) ->
    fs.readdir file, (err, files) ->
        next(err) if err
        movies = files.filter (file) ->
            file[0] isnt '.'
        res.json
            'movies' : movies

exports.play = (req, res, next) ->
    moviePath = path.join file, req.body.movie
    fs.readdir moviePath, (err, files) ->
        next(err) if err
        movieFiles = files.filter (file) ->
            path.extname(file) in ['.mp4', '.avi', '.mov']
        if movieFiles.length > 0
            # omx.play path.join(moviePath, movieFiles[0])
            console.log 'Playing movie at ' + path.join(moviePath, movieFiles[0])
        else
            next(new Error 'Movie file not found')
        res.send 200