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
            res.json
                movies : results.sort()

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

# exports.test = (req, res, next) ->
#     movie = '/Users/Joe/Movies/2012 NBA TV\'s The Dream Team 720p.mp4'
#     stats = fs.statSync movie
#     range = req.range stats.size
#     console.log range[0]
#     stream = fs.createReadStream movie,
#         start : range[0].start
#         end : range[0].end
#     res.status 200
#     res.type '.mp4'
#     res.set
#         'Accept-Ranges' : 'bytes'
#         'Content-Range' : "bytes #{range[0].start}-#{range[0].end}/#{stats.size}"
#         'Content-Length' : stats.size
#     stream.pipe res