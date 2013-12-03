fs     = require 'fs'
path   = require 'path'
config = require './config'
omx    = require './lib/omxcontrol'

exports.get = (req, res, next) ->
    fs.readdir config.youtubeDir, (err, videos) ->
        next(err) if err?
        titles = videos.map (video) ->
            filename : video
            title : video.substr 0, video.lastIndexOf('-')
        res.json
            'videos' : titles

exports.play = (req, res, next) ->
    video = path.join config.youtubeDir, req.body.filename
    fs.exists video, (exists) ->
        if exists
            omx.quit()
            omx.start video
        else
            next(new Error('YouTube video could not be found.'))
        res.send 200