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

services.factory 'Movies', ['$resource', '$rootScope', 'Player', '$cacheFactory', ($resource, $rootScope, Player, $cacheFactory) ->
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
        Movies.play({id : id}).$promise.then () ->
            Player.isPlaying true
            Player.nowPlaying {movie : id}
            $rootScope.$broadcast 'play'

    api.delete = (id, deleteFile) ->
        Movies.delete({id : id, file : deleteFile}).$promise

    api.scan = () ->
        Movies.scan().$promise

    api.clearCache = () ->
        movieCache.removeAll()

    return api
]

services.factory 'Shows', ['$resource', '$rootScope', 'Player', '$route', '$cacheFactory', ($resource, $rootScope, Player, $route, $cacheFactory) ->
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
        Shows.playEpisode({id : id}).$promise.then () ->
            Player.isPlaying true
            Player.nowPlaying
                episode : id
                show : $route.current.params.id
                season : $route.current.params.season
            $rootScope.$broadcast 'play'

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

services.factory 'Player', ['$rootScope', 'playerCommands', '$resource', ($rootScope, playerCommands, $resource) ->
    Player = $resource '/player/:command', {command : '@command'}

    api = {}
    api.isPaused = (isPaused) ->
        if isPaused?
            localStorage['isPaused'] = isPaused
        else
            return if localStorage['isPaused'] is 'true' then true else false
    api.isPlaying = (isPlaying) ->
        if isPlaying?
            localStorage['isPlaying'] = isPlaying
        else
            return if localStorage['isPlaying'] is 'true' then true else false
    api.nowPlaying = (nowPlaying) ->
        if nowPlaying? and angular.isObject(nowPlaying)
            localStorage['nowPlaying'] = JSON.stringify nowPlaying
        else
            playing = localStorage['nowPlaying']
            return if playing? then JSON.parse(playing) else playing
    api.toggle = () ->
        Player.get {command : playerCommands.TOGGLE}
        api.isPaused(not api.isPaused())
    api.backward = () ->
        Player.get {command : playerCommands.BACKWARD}
    api.forward = () ->
        Player.get {command : playerCommands.FORWARD}
    api.stop = () ->
        Player.get {command : playerCommands.STOP}
        localStorage.clear()
        $rootScope.$broadcast 'stop'
    api.fastBackward = () ->
        Player.get {command : playerCommands.FASTBACKWARD}
    api.fastForward = () ->
        Player.get {command : playerCommands.FASTFORWARD}
    api.clearCache = () ->
        localStorage.clear()
        $rootScope.$broadcast 'stop'

    return api
]
