raspTv = angular.module 'raspTv', ['ngRoute', 'ngAnimate', 'angular-loading-bar', 'ngResource', 'mgcrea.ngStrap', 'simongeeks.filters', 'raspTv.services']

raspTv.config ['$routeProvider', '$httpProvider', ($routeProvider, $httpProvider) ->
    $httpProvider.interceptors.push 'errorInterceptor'

    $routeProvider.when '/movies',
        templateUrl : '/templates/movies.html'
        controller : 'movieCtrl'
        resolve :
            movies : ['Movies', (Movies) ->
                Movies.getAll true
            ]
    $routeProvider.when '/:type/:id/play',
        templateUrl : '/templates/play.html'
        controller : 'playCtrl'
        resolve :
            playing : ['$route', 'Shows', 'Movies', ($route, Shows, Movies) ->
                if $route.current.params.type is 'movies'
                    Movies.get $route.current.params.id
                else
                    Shows.getEpisode($route.current.params.id).then (episode) ->
                        Shows.get(episode.showId).then (show) ->
                            title: show.title
                            episode: episode
            ]
    $routeProvider.when '/:type/:id/mode',
        templateUrl : '/templates/mode.html'
        controller : 'modeCtrl'
        resolve :
            title : ['$route', 'Movies', 'Shows', ($route, Movies, Shows) ->
                if $route.current.params.type is 'movies'
                    Movies.get($route.current.params.id).then (movie) -> movie.title
                else
                    Shows.getEpisode($route.current.params.id).then (episode) ->
                        Shows.get(episode.showId).then (show) ->
                            "#{show.title} - #{episode.season} - #{episode.title}"
            ]
    $routeProvider.when '/:type/:id/stream',
        templateUrl : '/templates/stream.html'
        controller : 'streamCtrl'
        resolve :
            title : ['$route', 'Movies', 'Shows', ($route, Movies, Shows) ->
                if $route.current.params.type is 'movies'
                    Movies.get($route.current.params.id).then (movie) -> movie.title
                else
                    Shows.getEpisode($route.current.params.id).then (episode) ->
                        Shows.get(episode.showId).then (show) ->
                            "#{show.title} - #{episode.season} - #{episode.title}"
            ]
    $routeProvider.when '/shows',
        templateUrl : '/templates/shows.html'
        controller : 'showsCtrl'
        resolve :
            shows : ['Shows', (Shows) ->
                Shows.getAll()
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
    $routeProvider.when '/edit',
        templateUrl : '/templates/edit.html'
        controller : 'editCtrl'
        resolve:
            nonIndexedMovies : ['Movies', (Movies) ->
                Movies.clearCache()
                Movies.getAll false
            ]
            nonIndexedEpisodes : ['Shows', (Shows) ->
                Shows.clearCache()
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
                Shows.getAll()
            ]
    $routeProvider.otherwise
        redirectTo : '/movies'
]

raspTv.run ['$rootScope', 'Player', ($rootScope, Player) ->
    getSession = () ->
        Player.getSession().then (session) ->
            $rootScope.session = if session.movieId? or session.episodeId? then session else null

    # get session initially
    getSession().then () ->
        # watch session and update the server when it changes or clear it if set to null
        $rootScope.$watch 'session', ((newVal, oldVal) ->
            if angular.equals(newVal, oldVal) then return
            if newVal? then Player.updateSession(newVal) else Player.clearSession()
        ), true

    $rootScope.refreshSession = getSession

    # Update session when moving to different pages.
    $rootScope.$on '$routeChangeSuccess', getSession
]

raspTv.controller 'navCtrl', ['$scope', '$location', '$rootScope', ($scope, $location, $rootScope) ->

    $scope.isActive = (page) ->
        if page is 'movies' and /^\/movies/.test($location.path()) and not /play$/.test($location.path())
            true
        else if page is 'shows' and (/^\/shows/.test($location.path()) or /^\/episodes/.test($location.path())) and not /play$/.test($location.path())
            true
        else if page is 'play' and /play$/.test($location.path())
            true
        else if page is 'edit' and /^\/edit/.test($location.path())
            true
        else
            false

    setUpLink = () ->
        if not $scope.session? then return

        if $scope.session.movieId?
            $scope.nowPlayingLink = "#/movies/#{$scope.session.movieId}/play"
        else if $scope.session.episodeId?
            $scope.nowPlayingLink = "#/episodes/#{$scope.session.episodeId}/play"

    # update now playing link if the session changes
    $rootScope.$watch 'session', setUpLink, true
]

raspTv.controller 'alertCtrl', ['$scope', '$rootScope', ($scope, $rootScope) ->

    $rootScope.$on 'alert', (event, alert) ->
        $scope.alert = alert

    $scope.close = () -> $scope.alert = null
]

raspTv.controller 'streamCtrl', ['$scope', '$routeParams', 'title', ($scope, $routeParams, title) ->
    $scope.title = title
    if $routeParams.type is 'movies'
        $scope.src = "/movies/#{$routeParams.id}/stream"
    else
        $scope.src = "/shows/episodes/#{$routeParams.id}/stream"
]

raspTv.controller 'movieCtrl', ['$scope', 'movies', 'Movies', ($scope, movies, Movies) ->
    $scope.movies = movies
    $scope.activePanel = -1

    $scope.scan = () ->
        Movies.scan().then () ->
            $scope.$emit 'alert',
                type : 'success'
                title : 'Success!'
]

raspTv.controller 'playCtrl', ['$scope', 'Player', 'Shows', 'Movies', '$routeParams', '$location', 'playing', '$rootScope', ($scope, Player, Shows, Movies, $routeParams, $location, playing, $rootScope) ->
    $scope.playing = playing
    if $routeParams.type is 'movies'
        $scope.isShow = false
        returnPath = '/movies'
    else if $routeParams.type is 'episodes'
        $scope.isShow = true
        returnPath = '/shows'
    else
        $location.path '/'

    setup = () ->
        $scope.toggle = () -> Player.toggle().then () -> $scope.session.isPaused = not $scope.session.isPaused
        $scope.backward = Player.backward
        $scope.forward = Player.forward
        $scope.fastBackward = Player.fastBackward
        $scope.fastForward = Player.fastForward
        $scope.stop = () ->
            Player.stop().finally () ->
                $rootScope.session = null
                $location.path returnPath

    if $scope.session? then setup()
    else
        promise = if $scope.isShow then Shows.play($routeParams.id) else Movies.play($routeParams.id)
        promise.then(setup).then $scope.refreshSession, () ->
            $location.path returnPath
]

raspTv.controller 'modeCtrl', ['$scope', '$routeParams', 'title', ($scope, $routeParams, title) ->
    $scope.title = title
    $scope.href  = "#/#{$routeParams.type}/#{$routeParams.id}"
]

raspTv.controller 'showsCtrl', ['$scope', 'shows', 'Shows', ($scope, shows, Shows) ->
    $scope.shows = shows

    $scope.scan = () ->
        Shows.scan().then () ->
            $scope.$emit 'alert',
                type : 'success'
                title : 'Success!'
]

raspTv.controller 'seasonsCtrl', ['$scope', 'show', '$location', ($scope, show, $location) ->
    $scope.show = show

    $scope.seasons = []
    for e in show.episodes when $scope.seasons.indexOf(e.season) is -1
        $scope.seasons.push e.season

    $scope.seasons = $scope.seasons.sort (a, b) ->
        parseInt(a, 10) - parseInt(b, 10)

    $scope.random = () ->
        season = $scope.seasons[Math.floor(Math.random() * $scope.seasons.length)]
        episodes = show.episodes.filter (e) -> e.season is season
        episodeId = episodes[Math.floor(Math.random() * episodes.length)].id
        $location.path "/episodes/#{episodeId}/mode"
]

raspTv.controller 'episodesCtrl', ['$scope', 'show', '$routeParams', '$location', ($scope, show, $routeParams, $location) ->
    $scope.show     = show
    $scope.season   = parseInt $routeParams.season, 10
    $scope.episodes = show.episodes.filter((e) -> e.season is $scope.season).sort (a, b) -> a.number - b.number

    $scope.random = () ->
        episodeId = $scope.episodes[Math.floor(Math.random() * $scope.episodes.length)].id
        $location.path "/episodes/#{episodeId}/mode"
]

raspTv.controller 'editCtrl', ['$scope', 'nonIndexedMovies', 'nonIndexedEpisodes', ($scope, movies, episodes) ->
    $scope.movies   = movies
    $scope.episodes = episodes
]

raspTv.controller 'editMovieCtrl', ['$scope', 'movie', 'Movies', '$location', '$window', ($scope, movie, Movies, $location, $window) ->
    $scope.movie = movie

    if not $scope.movie.isIndexed
        $scope.movie.title = $scope.movie.filepath.substring($scope.movie.filepath.lastIndexOf('/') + 1, $scope.movie.filepath.lastIndexOf('.'))

    $scope.save = () ->
        Movies.save($scope.movie).then () ->
            $scope.$emit 'alert',
                type : 'success'
                title : 'Success!'
                msg : "#{$scope.movie.title} was updated."
            Movies.clearCache()
            if not $scope.movie.isIndexed then $location.path('/edit')


    $scope.deleteMovie = (deleteFile) ->
        if $window.confirm('Are you sure you want to delete this movie?')
            Movies.delete($scope.movie.id, deleteFile).then () ->
                $scope.$emit 'alert',
                        type : 'success'
                        title : 'Success!'
                        msg : "#{$scope.movie.title} was deleted."
                    Movies.clearCache()
                    if not $scope.movie.isIndexed then $location.path('/edit')

]
raspTv.controller 'editEpisodeCtrl', ['$scope', 'episode', 'shows', 'Shows', '$location', '$window', ($scope, episode, shows, Shows, $location, $window) ->
    $scope.episode = episode
    $scope.shows = shows

    # Set default title
    if not $scope.episode.isIndexed
        $scope.episode.title = $scope.episode.filepath.substring($scope.episode.filepath.lastIndexOf('/') + 1, $scope.episode.filepath.lastIndexOf('.'))

        # Try to get season and episode number from filename
        results = /[sS](0?\d+)[eE](0?\d+)/.exec($scope.episode.filepath)
        matches = if results? then results[1..] else []
        if matches[0]? then $scope.episode.season = parseInt matches[0], 10
        if matches[1]? then $scope.episode.number = parseInt matches[1], 10

    $scope.saveShow = () ->
        Shows.clearCache()
        Shows.add($scope.show).then(Shows.getAll).then (shows) ->
            $scope.shows = shows
            $scope.show = ''

    $scope.saveEpisode = () ->
        $scope.episode.number = parseInt $scope.episode.number, 10
        $scope.episode.season = parseInt $scope.episode.season, 10
        $scope.episode.showId = parseInt $scope.episode.showId, 10
        Shows.saveEpisode($scope.episode).then () ->
            $scope.$emit 'alert',
                    type : 'success'
                    title : 'Success!'
                    msg : "#{$scope.episode.title} was updated."
                Shows.clearCache()
                if not $scope.episode.isIndexed then $location.path('/edit')

    $scope.deleteEpisode = (deleteFile) ->
        if $window.confirm('Are you sure you want to delete this episode?')
            Shows.deleteEpisode($scope.episode.id, deleteFile).then () ->
                $scope.$emit 'alert',
                        type : 'success'
                        title : 'Success!'
                        msg : "#{$scope.episode.title} was deleted."
                Shows.clearCache()
                if not $scope.episode.isIndexed then $location.path('/edit')
]
