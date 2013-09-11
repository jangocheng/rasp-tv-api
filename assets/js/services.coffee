services = angular.module 'raspTv.services', ['LocalStorageModule', 'btford.socket-io']

services.factory 'movies', ['$http', ($http) ->
    movies = []
    movie = ''
    api = {}

    api.getAll = (cb) ->
        if movies.length is 0
            req = $http.get '/movies'
            req.success (data) ->
                movies = data.movies
                cb null, movies
            req.error (err) ->
                cb err, null
        else
            cb null, movies
    return api
]

services.factory 'shows', ['$http', ($http) ->
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

    api.getAll = (cb) ->
        if shows.length is 0
            req = $http.get '/shows'
            req.success (data) ->
                shows = data.shows
                cb null, shows
            req.error (err) ->
                cb err, null
        else
            cb null, shows

    api.getSeasons = (cb) ->
        req = $http.post '/shows/seasons',
            'show' : show
        req.success (data) ->
            cb null, data.seasons
        req.error (err) ->
            cb err, null

    api.getEpisodes = (cb) ->
        req = $http.post '/shows/seasons/episodes',
            'show' : show
            'season' : season
        req.success (data) ->
            episodes = data.episodes.map (item) ->
                results = /^(\d+)\s-\s(.+)\.\w+$/.exec item
                {filename : item, name : results[2], num : results[1]}
            cb null, episodes
        req.error (err) ->
            cb err, null

    api.getEpisode = getEpisode
    api.setEpisode = setEpisode
    api.getSeason = getSeason
    api.setSeason = setSeason
    api.getShow = getShow
    api.setShow = setShow

    return api
]

services.factory 'player', ['$rootScope', '$http', 'localStorageService', 'socket', ($rootScope, $http, localStorageService, socket) ->
    nowPlaying = ''
    isPaused   = false
    isPlaying  = false

    # Getters and setters
    setNowPlaying = (video) ->
        nowPlaying = video
        localStorageService.set 'playing', video
    getNowPlaying = () ->
        nowPlaying
    getIsPaused = () ->
        isPaused
    setIsPaused = (paused) ->
        isPaused = paused
        localStorageService.set 'isPaused', paused
    getIsPlaying = () ->
        isPlaying
    setIsPlaying = (playing) ->
        isPlaying = playing
        localStorageService.set 'isPlaying', playing
        $rootScope.isPlaying = playing

    api = {}
    api.checkCache = () ->
        nowPlaying = localStorageService.get 'playing'
        isPaused   = if localStorageService.get('isPaused') is 'true' then true else false
        isPlaying  = if localStorageService.get('isPlaying') is 'true' then true else false
        $rootScope.isPlaying = isPlaying

    api.playMovie = (movie, cb) ->
        req = $http.post '/movies/play',
            'movie' : movie
        req.success () ->
            setNowPlaying movie
            setIsPlaying true
            setIsPaused false
            cb()
        req.error (err) ->
            cb err

    api.playShow = (show, season, episode, cb) ->
        showProps =
            'show' : show
            'season' : season
            'episode' : episode.filename
        req = $http.post '/shows/play', showProps
        req.success () ->
            showProps.episode = episode.name
            setNowPlaying showProps
            setIsPlaying true
            setIsPaused false
            cb()
        req.error (err) ->
            cb err

    api.playRandomEpisode = (show, cb) ->
        req = $http.post '/shows/random', {'show' : show}
        req.success (data) ->
            showProps = data
            showProps.show = show
            matches = /^\d+\s-\s(.+)\.\w+$/.exec showProps.episode
            showProps.episode = matches[1]
            setNowPlaying showProps
            setIsPlaying true
            setIsPaused false
            cb()
        req.error (err) ->
            cb err

    api.toggle = () ->
        socket.emit 'toggle'
        setIsPaused not isPaused
    api.backward = () ->
        socket.emit 'backward'
    api.forward = () ->
        socket.emit 'forward'
    api.stop = () ->
        socket.emit 'stop'
        setIsPlaying false
        localStorageService.clearAll()

    api.setNowPlaying = setNowPlaying
    api.getNowPlaying = getNowPlaying
    api.isPaused = getIsPaused
    api.setIsPaused = setIsPaused
    api.isPlaying = getIsPlaying
    api.setIsPlaying = setIsPlaying

    return api
]
