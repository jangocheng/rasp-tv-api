services = angular.module 'raspTv.services', []

services.constant 'playerCommands',
    TOGGLE : 0
    BACKWARD : 1
    FORWARD : 2
    STOP : 3
    FASTBACKWARD : 4
    FASTFORWARD : 5

services.factory 'errorInterceptor', ['$q', '$rootScope', ($q, $rootScope) ->

    broadcastError = (err) ->
        $rootScope.$broadcast 'alert',
            type : 'error'
            title  : "Error: #{err.config.method} #{err.status} #{err.config.url}"
            msg : err.data.error
        $q.reject err
    {
        'requestError'  : broadcastError
        'responseError' : broadcastError
    }
]

services.factory 'Movies', ['$resource', '$cacheFactory', ($resource, $cacheFactory) ->
    api = {}
    movieCache = $cacheFactory 'movie'
    Movies = $resource '/movies/:id', {id : '@id'},
        query :
            method : 'GET'
            cache : movieCache
            isArray : true
        get :
            method : 'GET'
            cache : movieCache
        play :
            method : 'GET'
            url : '/movies/:id/play'
            params : {id : '@id'}
        scan :
            method : 'GET'
            url : '/scan/movies'

    api.getAll = (isIndexed) ->
        Movies.query({isIndexed : isIndexed}).$promise

    api.get = (id) ->
        Movies.get({id : id}).$promise

    api.save = (movie) ->
        Movies.save({id : movie.Id}, movie).$promise

    api.play = (id) ->
        Movies.play({id : id}).$promise

    api.delete = (id, deleteFile) ->
        Movies.delete({id : id, file : deleteFile}).$promise

    api.scan = () ->
        Movies.scan().$promise

    api.clearCache = () ->
        movieCache.removeAll()

    return api
]

services.factory 'Shows', ['$resource', '$cacheFactory', ($resource, $cacheFactory) ->
    api = {}
    showsCache = $cacheFactory 'shows'
    Shows = $resource '/shows/:id', {id : '@id'},
        query :
            method : 'GET'
            isArray : true
            cache : showsCache
        get :
            method : 'GET'
            cache : showsCache
        add :
            method : 'POST'
            url : '/shows/add'
        getEpisode :
            method : 'GET'
            url : '/shows/episodes/:id'
            params : {id : '@id'}
            cache : showsCache
        getAllEpisodes :
            method : 'GET'
            url : '/episodes'
            isArray : true
            cache : showsCache
        saveEpisode :
            method : 'POST'
            url : '/shows/episodes/:id'
            params : {id : '@id'}
        playEpisode :
            method : 'GET'
            url : '/shows/episodes/:id/play'
            params : {id : '@id'}
        scan :
            method : 'GET'
            url : '/scan/episodes'
        deleteEpisode :
            method: 'DELETE'
            url : '/shows/episodes/:id'
            params : {id : '@id'}

    api.getAll = () ->
        Shows.query().$promise

    api.get = (id) ->
        Shows.get({id : id}).$promise

    api.add = (show) ->
        Shows.add({}, show).$promise

    api.getAllEpisodes = (isIndexed) ->
        Shows.getAllEpisodes({isIndexed : isIndexed}).$promise

    api.getEpisode = (id) ->
        Shows.getEpisode({id : id}).$promise

    api.saveEpisode = (episode) ->
        Shows.saveEpisode({id : episode.Id}, episode).$promise

    api.play = (id) ->
        Shows.playEpisode({id : id}).$promise

    api.scan = () ->
        Shows.scan().$promise

    api.getEpisodeFromShow = (show, episodeId) ->
        for e in show.Episodes
            if e.Id is episodeId
                return e

    api.deleteEpisode = (id, deleteFile) ->
        Shows.deleteEpisode({id : id, file : deleteFile}).$promise

    api.clearCache = () ->
        showsCache.removeAll()

    return api
]

services.factory 'Player', ['playerCommands', '$resource', (playerCommands, $resource) ->
    Player  = $resource '/player/command/:command', {command : '@command'}
    Session = $resource '/player/session'

    api = {}
    api.getSession = () ->
        Session.get().$promise
    api.clearSession = () ->
        Session.remove().$promise
    api.updateSession = (session) ->
        Session.save(session).$promise
    api.toggle = () ->
        Player.get({command : playerCommands.TOGGLE}).$promise
    api.backward = () ->
        Player.get({command : playerCommands.BACKWARD}).$promise
    api.forward = () ->
        Player.get({command : playerCommands.FORWARD}).$promise
    api.stop = () ->
        Player.get({command : playerCommands.STOP}).$promise
    api.fastBackward = () ->
        Player.get({command : playerCommands.FASTBACKWARD}).$promise
    api.fastForward = () ->
        Player.get({command : playerCommands.FASTFORWARD}).$promise

    return api
]
