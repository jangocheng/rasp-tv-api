fs    = require 'fs'
path  = require 'path'
async = require 'async'
omx   = require './lib/omxcontrol'

showsDir  = '/media/passport/TV Shows'
# showsDir = '/Volumes/My Passport/TV Shows'

exports.index = (req, res, next) ->
    fs.readdir showsDir, (err, files) ->
        next(err) if err?
        async.filter files, ((file, cb) ->
            if file[0] is '.'
                cb false
            else
                fs.stat path.join(showsDir, file), (err, stats) ->
                    cb(false) if err?
                    cb stats.isDirectory()
        ), (results) ->
            res.json
                shows : results.sort()

exports.seasons = (req, res, next) ->
    showPath = path.join showsDir, req.body.show
    fs.readdir showPath, (err, files) ->
        next(err) if err?
        async.filter files, ((file, cb) ->
            if file[0] is '.'
                cb false
            else
                fs.stat path.join(showPath, file), (err, stats) ->
                    cb(false) if err?
                    cb stats.isDirectory()
        ), (results) ->
            res.json
                seasons : results.sort()

exports.episodes = (req, res, next) ->
    showPath = path.join showsDir, req.body.show, '' + req.body.season
    fs.readdir showPath, (err, files) ->
        next(err) if err?
        episodes = files.filter (file) ->
            file[0] isnt '.' and path.extname(file) in ['.mp4', '.avi', '.mov', '.mkv']
        res.json
            'episodes' : episodes.sort (a, b) ->
                test = /^(\d+)\s-\s.+/
                matchesA = test.exec a
                matchesB = test.exec b
                episodeNumA = parseInt matchesA[1], 10
                episodeNumB = parseInt matchesB[1], 10
                return 1 if episodeNumA > episodeNumB
                return -1 if episodeNumA < episodeNumB
                return 0

exports.play = (req, res, next) ->
    episodePath = path.join showsDir, req.body.show, '' + req.body.season, req.body.episode
    fs.exists episodePath, (exists) ->
        if exists
            omx.quit()
            omx.start episodePath
        else
            next(new Error 'TV Show episode not found')
        res.send 200
