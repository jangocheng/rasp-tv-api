module.exports = (grunt) ->
    grunt.initConfig
        pkg : grunt.file.readJSON 'package.json'
        coffee :
            'assets/js/app.js' : 'assets/js/app.coffee'
        clean : ['app.js', 'routes/*.js', 'assets/js/*.js', 'assets/templates/*.html']
        watch :
            coffee :
                files : ['assets/js/*.coffee']
                tasks : 'coffee'
            jade :
                files : ['assets/templates/*.jade']
                tasks : 'jade'
        jade :
            all :
                files :
                    'assets/templates/movies.html' : 'assets/templates/movies.jade'

    grunt.loadNpmTasks 'grunt-contrib-coffee'
    grunt.loadNpmTasks 'grunt-contrib-clean'
    grunt.loadNpmTasks 'grunt-contrib-watch'
    grunt.loadNpmTasks 'grunt-contrib-jade'
    grunt.registerTask 'default', ['coffee', 'jade']