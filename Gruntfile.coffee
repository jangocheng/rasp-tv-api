module.exports = (grunt) ->
    grunt.initConfig
        pkg : grunt.file.readJSON 'package.json'
        coffee :
            'app.js' : 'app.coffee'
            'routes/*.js' : 'routes/*.coffee'
            'assets/js/*.js' : 'assets/js/*.coffee'
        clean : ['app.js', 'routes/*.js', 'assets/js/*.js']
        watch :
            coffee :
                files : ['app.coffee', 'routes/*.coffee', 'assets/js/*.coffee']
                tasks : 'coffee'

    grunt.loadNpmTasks 'grunt-contrib-coffee'
    grunt.loadNpmTasks 'grunt-contrib-clean'
    grunt.loadNpmTasks 'grunt-contrib-watch'
    grunt.registerTask 'default', 'coffee'