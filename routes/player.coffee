omx     = require './lib/omxcontrol'
path    = require 'path'
config  = require './config'
youtube = require 'youtube-dl'

TOGGLE = 0
BACKWARD = 1
FORWARD = 2
STOP = 3
FASTBACKWARD = 4
FASTFORWARD = 5

module.exports = (socket) ->
    socket.on 'message', (data, flags) ->
        if isNaN data
            msg = JSON.parse data
            download = youtube.download msg.url, config.youtubeDir
            download.on 'progress', (data) ->
                socket.send JSON.stringify
                    percent : data.percent
            download.on 'end', (info) ->
                if msg.shouldPlay
                    omx.quit()
                    omx.start path.join(config.youtubeDir, info.filename)
                socket.send JSON.stringify
                    end :
                        title : info.filename.substr 0, info.filename.indexOf(info.id) - 1
            download.on 'error', (err) ->
                # socket.emit 'error', err
        else
            command = parseInt data, 10
            switch command
                when TOGGLE
                    omx.pause()
                when BACKWARD
                    omx.backward()
                when FORWARD
                    omx.forward()
                when STOP
                    omx.quit()
                when FASTBACKWARD
                    omx.fastBackward()
                when FASTFORWARD
                    omx.fastForward()
