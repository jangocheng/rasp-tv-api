omx     = require './lib/omxcontrol'
path    = require 'path'
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
        youtubeDir = '/media/passport/youtube'
        youtube.getInfo data.url, (err, info) ->
            throw err if err?
            download = youtube.download data.url, youtubeDir
            download.on 'progress', (data) ->
                socket.emit 'progress', {percent : data.percent}
            download.on 'end', (data) ->
                omx.quit()
                omx.start path.join(youtubeDir, data.filename)
                socket.emit 'end',
                    title : info.title
            download.on 'error', (err) ->
                socket.emit 'error', err