services = angular.module 'raspTv.services', []

services.constant 'playerCommands',
    TOGGLE : 0
    BACKWARD : 1
    FORWARD : 2
    STOP : 3
    FASTBACKWARD : 4
    FASTFORWARD : 5

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

services.factory 'player', ['$rootScope', '$http', '$location', 'playerCommands', ($rootScope, $http, $location, playerCommands) ->
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

    api.playYoutube = (video, cb) ->
        req = $http.post '/youtube/play', {'filename' : video.filename}
        req.success () ->
            setNowPlaying video.title
            setIsPlaying true
            setIsPaused false
            cb()
        req.error (err) ->
            cb err

    api.downloadYoutube = (url, shouldPlay, cb) ->
        downloadProgress = 0
        socket.send JSON.stringify({'url', url, 'shouldPlay' : shouldPlay})
        socket.onmessage = (event) ->
            $rootScope.$apply () ->
                data = JSON.parse event.data
                if data.percent?
                    downloadProgress = data.percent
                    $rootScope.$broadcast 'progress', downloadProgress
                else if data.end?
                    downloadProgress = 0
                    if shouldPlay
                        setNowPlaying data.end.title
                        setIsPlaying true
                        setIsPaused false
                    cb()
        socket.onerror = (err) ->
            cb err

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

services.factory 'youtube', ['$http', ($http) ->

    videos = []
    api = {}

    api.getAll = (cb) ->
        if videos.length is 0
            req = $http.get '/youtube/videos'
            req.success (data) ->
                videos
                cb null, data.videos
            req.error (err) ->
                cb err, null
        else
            cb null, videos

    return api
]
