omx = require './lib/omxcontrol'

module.exports = (socket) ->
    socket.on 'pause', (data) ->
        if data.isPaused then omx.play() else omx.pause()

    socket.on 'play', (data) ->
        omx.play()

    socket.on 'backward', (data) ->
        omx.backward()

    socket.on 'forward', (data) ->
        omx.forward()