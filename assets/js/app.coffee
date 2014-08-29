raspTv = angular.module 'raspTv', ['ngRoute', 'ngAnimate', 'angular-loading-bar', 'ngResource', 'raspTv.services']

raspTv.config ['$routeProvider', '$httpProvider', ($routeProvider, $httpProvider) ->
    $httpProvider.interceptors.push 'errorInterceptor'

    $routeProvider.when '/movies',
        templateUrl : '/templates/movies.html'
        controller : 'movieCtrl'
        resolve :
            movies : ['Movies', (Movies) ->
                Movies.getAll true
            ]
    $routeProvider.when '/movies/:id/play',
        templateUrl : '/templates/play.html'
        controller : 'playMovieCtrl'
        resolve :
            movie : ['$route', 'Movies', ($route, Movies) ->
                Movies.get $route.current.params.id
            ]
    $routeProvider.when '/movies/:id/mode',
        templateUrl : '/templates/mode.html'
        controller : 'modeCtrl'
    $routeProvider.when '/movies/:id/stream',
        templateUrl : '/templates/stream.html'
        controller : 'streamCtrl'
    $routeProvider.when '/shows',
        templateUrl : '/templates/shows.html'
        controller : 'showsCtrl'
        resolve :
            shows : ['Shows', (Shows) ->
                Shows.getAll true
            ]
    $routeProvider.when '/shows/:id/seasons',
        templateUrl : '/templates/seasons.html'
        controller : 'seasonsCtrl'
        resolve :
            show : ['$route', 'Shows', ($route, Shows) ->
                Shows.get $route.current.params.id
            ]
    $routeProvider.when '/shows/:id/seasons/:season/episodes',
        templateUrl : '/templates/episodes.html'
        controller : 'episodesCtrl'
        resolve :
            show : ['$route', 'Shows', ($route, Shows) ->
                Shows.get $route.current.params.id
            ]
    $routeProvider.when '/shows/:id/seasons/:season/episodes/:episode/play',
        templateUrl : '/templates/play.html'
        controller : 'playShowCtrl'
        resolve :
            show : ['$route', 'Shows', ($route, Shows) ->
                Shows.get $route.current.params.id
            ]
    $routeProvider.when '/shows/:id/seasons/:season/episodes/:episode/mode',
        templateUrl : '/templates/mode.html'
        controller : 'modeCtrl'
    $routeProvider.when '/shows/:id/seasons/:season/episodes/:episode/stream',
        templateUrl : '/templates/stream.html'
        controller : 'streamCtrl'
    $routeProvider.when '/edit',
        templateUrl : '/templates/edit.html'
        controller : 'editCtrl'
        resolve:
            nonIndexedMovies : ['Movies', (Movies) ->
                Movies.getAll false
            ]
            nonIndexedEpisodes : ['Shows', (Shows) ->
                Shows.getAllEpisodes false
            ]
    $routeProvider.when '/edit/movie/:id',
        templateUrl : '/templates/editMovie.html'
        controller : 'editMovieCtrl'
        resolve :
            movie : ['$route', 'Movies', ($route, Movies) ->
                Movies.get $route.current.params.id
            ]
    $routeProvider.when '/edit/episode/:id',
        templateUrl : '/templates/editEpisode.html'
        controller : 'editEpisodeCtrl'
        resolve :
            episode : ['$route', 'Shows', ($route, Shows) ->
                Shows.getEpisode $route.current.params.id
            ]
            shows : ['Shows', (Shows) ->
                Shows.getAll false
            ]
    $routeProvider.otherwise
        redirectTo : '/movies'
]

raspTv.controller 'navCtrl', ['$scope', '$location', 'Player', ($scope, $location, Player) ->
    $scope.isPlaying = Player.isPlaying()

    $scope.isActive = (page) ->
        if page is 'movies' and /^\/movies/.test($location.path()) and not /play$/.test($location.path())
            true
        else if page is 'shows' and /^\/shows/.test($location.path()) and not /play$/.test($location.path())
            true
        else if page is 'play' and /play$/.test($location.path())
            true
        else if page is 'edit' and /^\/edit/.test($location.path())
            true
        else
            false

    setUpLink = () ->
        nowPlaying = Player.nowPlaying()
        if nowPlaying?
            if nowPlaying.movie?
                $scope.nowPlayingLink = "#/movies/#{nowPlaying.movie}/play"
            else
                $scope.nowPlayingLink = "#/shows/#{nowPlaying.show}/seasons/#{nowPlaying.season}/episodes/#{nowPlaying.episode}/play"

    $scope.$on 'play', () ->
        setUpLink()
        $scope.isPlaying = true
    $scope.$on 'stop', () -> $scope.isPlaying = false

    setUpLink()
]

raspTv.controller 'errorCtrl', ['$scope', '$rootScope', ($scope, $rootScope) ->

    $scope.$on 'httpError', (event, err) -> $scope.error = err

    $scope.close = () -> $scope.error = null
]

raspTv.controller 'errorCtrl', ['$scope', '$rootScope', ($scope, $rootScope) ->

    $rootScope.$on 'httpError', (event, err) ->
        $scope.error = err

    $scope.close = () ->
        $scope.error = null
]

raspTv.controller 'streamCtrl', ['$scope', '$routeParams', ($scope, $routeParams) ->
    if $routeParams.episode?
        $scope.src = "/shows/episodes/#{$routeParams.episode}/stream"
    else
        $scope.src = "/movies/#{$routeParams.id}/stream"
]

raspTv.controller 'movieCtrl', ['$scope', 'movies', 'Movies', ($scope, movies, Movies) ->
    $scope.movies = movies

    $scope.scan = () ->
        Movies.scan().then () -> $scope.success = true

    $scope.close = () -> $scope.success = false
]

raspTv.controller 'playMovieCtrl', ['$scope', 'Player', 'movie', 'Movies', '$routeParams', '$location', ($scope, Player, movie, Movies, $routeParams, $location) ->
    setup = () ->
        $scope.isShow = false
        $scope.isPaused = Player.isPaused()
        $scope.playing = movie

        $scope.toggle = () ->
            Player.toggle()
            $scope.isPaused = Player.isPaused()

        $scope.backward = Player.backward
        $scope.forward = Player.forward
        $scope.fastBackward = Player.fastBackward
        $scope.fastForward = Player.fastForward
        $scope.stop = () ->
            Player.stop()
            $location.path '/movies'

    nowPlaying = Player.nowPlaying()
    if nowPlaying? and nowPlaying.movie? and nowPlaying.movie is $routeParams.id
        setup()
    else
        Movies.play($routeParams.id).then setup, () ->
            Player.clearCache()
            $location.path '/movies'
]

raspTv.controller 'modeCtrl', ['$scope', '$routeParams', ($scope, $routeParams) ->
    if $routeParams.episode?
        $scope.href = "#/shows/#{$routeParams.id}/seasons/#{$routeParams.season}/episodes/#{$routeParams.episode}"
    else
        $scope.href = "#/movies/#{$routeParams.id}"
]

raspTv.controller 'showsCtrl', ['$scope', 'shows', '$rootScope', '$location', ($scope, shows, $rootScope, $location) ->
    shows.getAll().then (shows) ->
        $scope.shows = shows
]

raspTv.controller 'playShowCtrl', ['$scope', 'Player', 'show', '$routeParams', 'Shows', '$location', ($scope, Player, show, $routeParams, Shows, $location) ->
    setup = () ->
        $scope.isShow = true
        $scope.isPaused = Player.isPaused()

        $scope.playing = show
        episodeId = parseInt $routeParams.episode, 10
        for e in show.Episodes
            if e.Id is episodeId
                $scope.playing.episode = e
                break

        $scope.toggle = () ->
            Player.toggle()
            $scope.isPaused = Player.isPaused()

        $scope.backward = Player.backward
        $scope.forward = Player.forward
        $scope.fastBackward = Player.fastBackward
        $scope.fastForward = Player.fastForward
        $scope.stop = () ->
            Player.stop()
            $location.path '/shows'

    nowPlaying = Player.nowPlaying()
    if nowPlaying? and nowPlaying.episode? and nowPlaying.episode is $routeParams.episode
        setup()
    else
        Shows.play($routeParams.episode).then setup, () ->
            Player.clearCache()
            $location.path '/shows'
]

raspTv.controller 'showsCtrl', ['$scope', 'shows', 'Shows', ($scope, shows, Shows) ->
    $scope.shows = shows
    $scope.scan = () ->
        Shows.scan().then () -> $scope.success = true

    $scope.close = () -> $scope.success = false
]

raspTv.controller 'seasonsCtrl', ['$scope', 'show', '$location', ($scope, show, $location) ->
    $scope.show = show

    $scope.seasons = []
    for e in show.Episodes when $scope.seasons.indexOf(e.Season.Int64) is -1
        $scope.seasons.push e.Season.Int64

    $scope.seasons = $scope.seasons.sort()

    $scope.random = () ->
        season = $scope.seasons[Math.floor(Math.random() * $scope.seasons.length)]
        episodes = (e for e in show.Episodes when e.Season.Int64 is season)
        episodeId = episodes[Math.floor(Math.random() * episodes.length)].Id
        $location.path "/shows/#{show.Id}/seasons/#{season}/episodes/#{episodeId}/mode"
]

raspTv.controller 'episodesCtrl', ['$scope', 'show', '$routeParams', '$location', ($scope, show, $routeParams, $location) ->
    $scope.show = show
    $scope.season = parseInt $routeParams.season, 10
    $scope.episodes = (e for e in show.Episodes when e.Season.Int64 is $scope.season).sort (a, b) ->
        a.Number.Int64 - b.Number.Int64

    $scope.random = () ->
        episodeId = $scope.episodes[Math.floor(Math.random() * $scope.episodes.length)].Id
        $location.path "/shows/#{show.Id}/seasons/#{$scope.season}/episodes/#{episodeId}/mode"
]

raspTv.controller 'editCtrl', ['$scope', 'nonIndexedMovies', 'nonIndexedEpisodes', ($scope, movies, episodes) ->
    $scope.movies = movies
    $scope.episodes = episodes
]
raspTv.controller 'editMovieCtrl', ['$scope', 'movie', 'Movies', '$location', ($scope, movie, Movies, $location) ->
    $scope.movie = movie

    $scope.save = () ->
        $scope.movie.Title.Valid = true
        Movies.save($scope.movie).then () ->
            $location.path '/edit'
]
raspTv.controller 'editEpisodeCtrl', ['$scope', 'episode', 'shows', 'Shows', '$location', ($scope, episode, shows, Shows, $location) ->
    $scope.episode = episode
    $scope.shows = shows

    $scope.saveShow = () ->
        Shows.add($scope.show).then () ->
            Shows.getAll().then (shows) ->
                $scope.shows = shows
                $scope.show = ''

    $scope.saveEpisode = () ->
        $scope.episode.Title.Valid = true
        $scope.episode.Number.Int64 = parseInt $scope.episode.Number.Int64, 10
        $scope.episode.Number.Valid = true
        $scope.episode.Season.Int64 = parseInt $scope.episode.Season.Int64, 10
        $scope.episode.Season.Valid = true
        $scope.episode.ShowId.Int64 = parseInt $scope.episode.ShowId.Int64, 10
        $scope.episode.ShowId.Valid = true
        Shows.saveEpisode($scope.episode).then () ->
            $location.path '/edit'
]