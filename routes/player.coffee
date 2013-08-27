omx = require './lib/omxcontrol'

module.exports = (socket) ->
    socket.on 'toggle', (data) ->
        omx.pause()
        # console.log 'pause'

    socket.on 'backward', (data) ->
        omx.backward()
        # console.log 'back'

    socket.on 'forward', (data) ->
        omx.forward()
        # console.log 'forward'

    socket.on 'stop', (data) ->
        omx.quit()
        # console.log 'stop'
