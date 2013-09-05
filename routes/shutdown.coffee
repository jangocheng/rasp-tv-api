exec = require('child_process').exec

module.exports = (req, res, next) ->
    exec 'shutdown -h now', (err, stdout, stderr) ->
        next(err) if err
        res.send 200