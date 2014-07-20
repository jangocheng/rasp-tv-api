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
        $rootScope.$broadcast 'httpError', err
        $q.reject err

    {
        'requestError'  : broadcastError
        'responseError' : broadcastError
    }
]

services.factory 'movies', ['$resource', '$q', ($resource, $q) ->
    movies = []
    api = {}
    movieService = $resource '/movies'

    api.getAll = () ->
        deferred = $q.defer()
        if movies.length is 0
            movieService.query (res) ->
                movies = res
                deferred.resolve movies
        else
            deferred.resolve movies
        deferred.promise

    return api
]

services.factory 'shows', ['$http', '$resource', '$q', ($http, $resource, $q) ->
    shows   = []
    episode = ''
    season  = ''
    show    = ''
    api     = {}

    # Getters and setters
    getEpisode = () ->
        episode
    setEpisode = (val) ->
        episode = val
    getSeason = () ->
        season
    setSeason = (val) ->
        season = val
    getShow = () ->
        show
    setShow = (val) ->
        show = val

    api.getAll = () ->
        deferred = $q.defer()
        if shows.length is 0
            $resource('/shows').query (res) ->
                shows = res
                deferred.resolve shows
        else
            deferred.resolve shows
        deferred.promise

    api.getSeasons = () ->
        deferred = $q.defer()
        req = $http.post '/shows/seasons',
            'show' : show
        req.success (data) ->
            deferred.resolve data.seasons
        req.error (err) ->
            deferred.reject err
        deferred.promise

    api.getEpisodes = () ->
        deferred = $q.defer()
        req = $http.post '/shows/seasons/episodes',
            'show' : show
            'season' : season
        req.success (data) ->
            episodes = data.episodes.map (item) ->
                results = /^(\d+)\s-\s(.+)\.\w+$/.exec item
                {filename : item, name : results[2], num : results[1]}
            deferred.resolve episodes
        req.error (err) ->
            deferred.reject err
        deferred.promise

    api.getEpisode = getEpisode
    api.setEpisode = setEpisode
    api.getSeason = getSeason
    api.setSeason = setSeason
    api.getShow = getShow
    api.setShow = setShow

    return api
]

services.factory 'player', ['$rootScope', '$http', '$location', 'playerCommands', '$q', ($rootScope, $http, $location, playerCommands, $q) ->
    nowPlaying = ''
    isPaused   = false
    isPlaying  = false
    socket     = new WebSocket "ws://#{$location.host()}:#{$location.port()}"

    # Getters and setters
    setNowPlaying = (video) ->
        nowPlaying = video
        if angular.isObject(video)
            localStorage.setItem 'playing', JSON.stringify(video)
        else
            localStorage.setItem 'playing', video
    getNowPlaying = () ->
        nowPlaying
    getIsPaused = () ->
        isPaused
    setIsPaused = (paused) ->
        isPaused = paused
        localStorage.setItem 'isPaused', paused
    getIsPlaying = () ->
        isPlaying
    setIsPlaying = (playing) ->
        isPlaying = playing
        localStorage.setItem 'isPlaying', playing
        $rootScope.isPlaying = playing

    api = {}
    api.checkCache = () ->
        item = localStorage.getItem('playing')
        nowPlaying = if item? and item[0] is '{' then JSON.parse(item) else item
        isPaused   = if localStorage.getItem('isPaused') is 'true' then true else false
        isPlaying  = if localStorage.getItem('isPlaying') is 'true' then true else false
        $rootScope.isPlaying = isPlaying

    api.playMovie = (movie, cb) ->
        deferred = $q.defer()
        req = $http.post '/movies/play',
            'movie' : movie
        req.success () ->
            setNowPlaying movie
            setIsPlaying true
            setIsPaused false
            deferred.resolve()
        req.error (err) ->
            deferred.reject err
        deferred.promise

    api.playShow = (show, season, episode, cb) ->
        showProps =
            'show' : show
            'season' : season
            'episode' : episode.filename
        deferred = $q.defer()
        req = $http.post '/shows/play', showProps
        req.success () ->
            showProps.episode = episode.name
            setNowPlaying showProps
            setIsPlaying true
            setIsPaused false
            deferred.resolve()
        req.error (err) ->
            deferred.reject err
        deferred.promise

    api.playRandomEpisode = (show, cb) ->
        deferred = $q.defer()
        req = $http.post '/shows/random', {'show' : show}
        req.success (data) ->
            showProps = data
            showProps.show = show
            matches = /^\d+\s-\s(.+)\.\w+$/.exec showProps.episode
            showProps.episode = matches[1]
            setNowPlaying showProps
            setIsPlaying true
            setIsPaused false
            deferred.resolve()
        req.error (err) ->
            deferred.reject err
        deferred.promise

    api.toggle = () ->
        socket.send playerCommands.TOGGLE
        setIsPaused not isPaused
    api.backward = () ->
        socket.send playerCommands.BACKWARD
    api.forward = () ->
        socket.send playerCommands.FORWARD
    api.stop = () ->
        socket.send playerCommands.STOP
        setIsPlaying false
        localStorage.clear()
    api.fastBackward = () ->
        socket.send playerCommands.FASTBACKWARD
    api.fastForward = () ->
        socket.send playerCommands.FASTFORWARD

    api.setNowPlaying = setNowPlaying
    api.getNowPlaying = getNowPlaying
    api.isPaused = getIsPaused
    api.setIsPaused = setIsPaused
    api.isPlaying = getIsPlaying
    api.setIsPlaying = setIsPlaying

    return api
]
