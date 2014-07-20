omx     = require './lib/omxcontrol'
path    = require 'path'
config  = require './config'

TOGGLE = 0
BACKWARD = 1
FORWARD = 2
STOP = 3
FASTBACKWARD = 4
FASTFORWARD = 5

module.exports = (socket) ->
    socket.on 'message', (data, flags) ->
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
