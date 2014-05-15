module.exports = (grunt) ->
    grunt.initConfig
        pkg : grunt.file.readJSON 'package.json'
        coffee :
            client :
                files :
                    'assets/js/app.js' : 'assets/js/app.coffee'
                    'assets/js/services.js' : 'assets/js/services.coffee'
            node :
                files :
                    'app.js' : 'app.coffee'
                    'routes/index.js' : 'routes/index.coffee'
                    'routes/movies.js' : 'routes/movies.coffee'
                    'routes/shows.js' : 'routes/shows.coffee'
                    'routes/player.js' : 'routes/player.coffee'
                    'routes/shutdown.js' : 'routes/shutdown.coffee'
                    'routes/youtube.js' : 'routes/youtube.coffee'
                    'routes/config.js' : 'routes/config.coffee'
        clean : ['app.js', 'routes/*.js', 'assets/js/*.js', 'assets/templates/*.html', 'dist', 'omxcontrol']
        watch :
            coffee :
                files : ['assets/js/*.coffee']
                tasks : 'coffee:client'
            jade :
                files : ['assets/templates/*.jade']
                tasks : 'jade'
        jade :
            client :
                files :
                    'assets/templates/movies.html' : 'assets/templates/movies.jade'
                    'assets/templates/play.html' : 'assets/templates/play.jade'
                    'assets/templates/shows.html' : 'assets/templates/shows.jade'
                    'assets/templates/seasons.html' : 'assets/templates/seasons.jade'
                    'assets/templates/episodes.html' : 'assets/templates/episodes.jade'
                    'assets/templates/youtube.html' : 'assets/templates/youtube.jade'
        rsync :
            options :
                args : ['--verbose']
                exclude : [
                    ".git",
                    "assets/templates/*.jade",
                    "node_modules",
                    ".DS_Store",
                    "Gruntfile.coffee",
                    ".gitignore",
                    "*.coffee",
                    "LICENSE",
                    "README.md",
                    "videos",
                    "omxcontrol"
                ]
                recursive : true
            dist :
                options :
                    src : './'
                    dest : './dist'
            production :
                options :
                    args : ['--verbose']
                    src : './dist/'
                    dest : '/home/joe/rasp-tv/'
                    host : 'joe@192.168.11.2'
                    syncDestIgnoreExcl: true
        shell :
            restart :
                command : "ssh joe@rpi 'sudo systemctl restart rasptv.service'"
        uglify :
            client :
                files :
                    'assets/js/rasptv.min.js' : [
                        'assets/js/libs/angular.js',
                        'assets/js/libs/angular-route.js',
                        'assets/js/services.js',
                        'assets/js/app.js'
                    ]

    grunt.loadNpmTasks 'grunt-contrib-coffee'
    grunt.loadNpmTasks 'grunt-contrib-clean'
    grunt.loadNpmTasks 'grunt-contrib-watch'
    grunt.loadNpmTasks 'grunt-contrib-jade'
    grunt.loadNpmTasks 'grunt-rsync'
    grunt.loadNpmTasks 'grunt-shell'
    grunt.loadNpmTasks 'grunt-contrib-uglify'
    grunt.registerTask 'default', ['coffee:client', 'jade']
    grunt.registerTask 'deploy', ['clean', 'coffee', 'jade', 'uglify', 'rsync', 'clean']