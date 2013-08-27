fs    = require 'fs'
path  = require 'path'
async = require 'async'
omx   = require './lib/omxcontrol'

moviesDir = '/media/passport/Movies'
# moviesDir = '/Volumes/My Passport/Movies'

exports.index = (req, res, next) ->
    fs.readdir moviesDir, (err, files) ->
        next(err) if err?
        async.filter files, ((file, cb) ->
            if file[0] is '.'
                cb false
            else
                fs.stat path.join(moviesDir, file), (err, stats) ->
                    cb(false) if err?
                    cb stats.isDirectory()
        ), (results) ->
            res.json
                movies : results.sort()

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