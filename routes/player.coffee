omx     = require './lib/omxcontrol'
path    = require 'path'
config  = require './config'
youtube = require 'youtube-dl'

module.exports = (socket) ->
    socket.on 'toggle', (data) ->
        omx.pause()

    socket.on 'backward', (data) ->
        omx.backward()

    socket.on 'forward', (data) ->
        omx.forward()

    socket.on 'stop', (data) ->
        omx.quit()

    socket.on 'fastBackward', (data) ->
        omx.fastBackward()

    socket.on 'fastForward', (data) ->
        omx.fastForward()

    socket.on 'youtube', (data) ->
        download = youtube.download data.url, config.youtubeDir
        download.on 'progress', (data) ->
            socket.emit 'progress', {percent : data.percent}
        download.on 'end', (data) ->
            omx.quit()
            omx.start path.join(config.youtubeDir, data.filename)
            socket.emit 'end',
                title : data.filename.substr 0, data.filename.indexOf(data.id) - 1
        download.on 'error', (err) ->
            socket.emit 'error', err