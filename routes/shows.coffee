fs    = require 'fs'
path  = require 'path'
async = require 'async'
omx   = require './lib/omxcontrol'

showsDir  = '/media/passport/TV Shows'
# showsDir = '/Volumes/My Passport/TV Shows'

getShows = (callback) ->
    fs.readdir showsDir, (err, files) ->
        callback(err, null) if err?
        async.filter files, ((file, cb) ->
            if file[0] is '.'
                cb false
            else
                fs.stat path.join(showsDir, file), (err, stats) ->
                    cb(false) if err?
                    cb stats.isDirectory()
        ), (results) ->
            callback null, results.sort()

getSeasons = (show, callback) ->
    showPath = path.join showsDir, show
    fs.readdir showPath, (err, files) ->
        callback(err, null) if err?
        async.filter files, ((file, cb) ->
            if file[0] is '.'
                cb false
            else
                fs.stat path.join(showPath, file), (err, stats) ->
                    cb(false) if err?
                    cb stats.isDirectory()
        ), (results) ->
            callback null, results.sort()

getEpisodes = (show, season, callback) ->
    showPath = path.join showsDir, show, '' + season
    fs.readdir showPath, (err, files) ->
        callback(err, null) if err?
        episodes = files.filter (file) ->
            file[0] isnt '.' and path.extname(file) in ['.mp4', '.avi', '.mov', '.mkv']
        callback null, episodes.sort (a, b) ->
            test = /^(\d+)\s-\s.+/
            matchesA = test.exec a
            matchesB = test.exec b
            episodeNumA = parseInt matchesA[1], 10
            episodeNumB = parseInt matchesB[1], 10
            return 1 if episodeNumA > episodeNumB
            return -1 if episodeNumA < episodeNumB
            return 0

play = (show, season, episode, callback) ->
    episodePath = path.join showsDir, show, '' + season, episode
    fs.exists episodePath, (exists) ->
        if exists
            omx.quit()
            omx.start episodePath
            callback()
        else
            callback(new Error 'TV Show episode not found')


exports.index = (req, res, next) ->
    getShows (err, shows) ->
        next(err) if err?
        res.json
            'shows' : shows

exports.seasons = (req, res, next) ->
    getSeasons req.body.show, (err, seasons) ->
        next(err) if err?
        res.json
            'seasons' : seasons

exports.episodes = (req, res, next) ->
    getEpisodes req.body.show, req.body.season, (err, episodes) ->
        next(err) if err?
        res.json
            'episodes' : episodes

exports.play = (req, res, next) ->
    play req.body.show, req.body.season, req.body.episode, (err) ->
        next(err) if err?
        res.send 200

exports.random = (req, res, next) ->
    getSeasons req.body.show, (err, seasons) ->
        next(err) if err?
        season = seasons[Math.floor(Math.random() * seasons.length)]
        getEpisodes req.body.show, season, (err, episodes) ->
            next(err) if err?
            episode = episodes[Math.floor(Math.random() * episodes.length)]
            play req.body.show, season, episode, (err) ->
                next(err) if err?
                res.json
                    'season' : season
                    'episode' : episode