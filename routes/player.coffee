omx = require './lib/omxcontrol'

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