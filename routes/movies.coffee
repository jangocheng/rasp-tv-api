fs     = require 'fs'
path   = require 'path'
async  = require 'async'
config = require './config'
omx    = require './lib/omxcontrol'

exports.index = (req, res, next) ->
    fs.readdir config.moviesDir, (err, files) ->
        next(err) if err?
        async.filter files, ((file, cb) ->
            if file[0] is '.'
                cb false
            else
                fs.stat path.join(config.moviesDir, file), (err, stats) ->
                    cb(false) if err?
                    cb stats.isDirectory()
        ), (results) ->
            res.json results.sort()

exports.play = (req, res, next) ->
    moviePath = path.join config.moviesDir, req.body.movie
    fs.readdir moviePath, (err, files) ->
        next(err) if err
        movieFiles = files.filter (file) ->
            path.extname(file) in config.supportedFormats
        if movieFiles.length > 0
            omx.quit()
            omx.start path.join(moviePath, movieFiles[0])
        else
            next(new Error 'Movie file not found')
        res.send 200