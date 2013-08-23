module.exports = (grunt) ->
    grunt.initConfig
        pkg : grunt.file.readJSON 'package.json'
        coffee :
            client :
                files :
                    'assets/js/app.js' : 'assets/js/app.coffee'
            node :
                files :
                    'app.js' : 'app.coffee'
                    'routes/index.js' : 'routes/index.coffee'
                    'routes/api.js' : 'routes/api.coffee'
                    'routes/player.js' : 'routes/player.coffee'
        clean : ['app.js', 'routes/*.js', 'assets/js/*.js', 'assets/templates/*.html', 'dist']
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
                    "README.md"
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
                    host : 'joe@rpi'
                    syncDestIgnoreExcl: true
        shell :
            deploy :
                command : "ssh joe@rpi 'cd rasp-tv && make restart'"

    grunt.loadNpmTasks 'grunt-contrib-coffee'
    grunt.loadNpmTasks 'grunt-contrib-clean'
    grunt.loadNpmTasks 'grunt-contrib-watch'
    grunt.loadNpmTasks 'grunt-contrib-jade'
    grunt.loadNpmTasks 'grunt-rsync'
    grunt.loadNpmTasks 'grunt-shell'
    grunt.registerTask 'default', ['coffee:client', 'jade']
    grunt.registerTask 'deploy', ['clean', 'coffee', 'jade', 'rsync', 'shell:deploy', 'clean']