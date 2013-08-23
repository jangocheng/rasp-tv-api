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
