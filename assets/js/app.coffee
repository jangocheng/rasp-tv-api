raspTv = angular.module 'raspTv', ['raspTv.services']

raspTv.config ['$routeProvider', ($routeProvider) ->
    $routeProvider.when '/',
        templateUrl : '/templates/movies.html'
        controller : 'movieCtrl'
    $routeProvider.when '/play',
        templateUrl : '/templates/play.html'
        controller : 'playCtrl'
    $routeProvider.when '/shows',
        templateUrl : '/templates/shows.html'
        controller : 'showsCtrl'
    $routeProvider.when '/shows/seasons',
        templateUrl : '/templates/seasons.html'
        controller : 'seasonsCtrl'
    $routeProvider.when '/shows/seasons/episodes',
        templateUrl : '/templates/episodes.html'
        controller : 'episodesCtrl'
]

raspTv.run ['player', (player) ->
    player.checkCache()
]

raspTv.controller 'navCtrl', ['$scope', '$location', ($scope, $location) ->
    $scope.isActive = (page) ->
        if page is 'movies' and $location.path() is '/'
            'active'
        else if page is 'shows' and /^\/shows.*/.test($location.path())
            'active'
        else if page is 'play' and $location.path() is '/play'
            'active'
        else
            ''
]

raspTv.controller 'errorCtrl', ['$scope', ($scope) ->
    $scope.close = () ->
        $scope.error = null
]

raspTv.controller 'movieCtrl', ['$scope', 'movies', '$rootScope', '$location', 'player', ($scope, movies, $rootScope, $location, player) ->
    movies.getAll (err, movies) ->
        if err?
            $rootScope.error = err.msg
        else
            $scope.movies = movies

    $scope.play = (index) ->
        player.playMovie $scope.movies[index], (err) ->
            if err?
                $rootScope.error = err.msg
            else
                $location.path 'play'
]

raspTv.controller 'playCtrl', ['$scope', '$location', 'player', '$rootScope', ($scope, $location, player, $rootScope) ->
    $scope.isPlaying = player.isPlaying()
    $location.path('/') if not $scope.isPlaying

    $scope.playing = player.getNowPlaying()
    $scope.isShow = angular.isObject $scope.playing
    $scope.isPaused = player.isPaused()

    $scope.toggle = () ->
        player.toggle()
        $scope.isPaused = player.isPaused()
    $scope.backward = player.backward
    $scope.forward = player.forward
    $scope.stop = () ->
        player.stop()
        $location.path '/'
]

raspTv.controller 'showsCtrl', ['$scope', 'shows', '$rootScope', '$location', ($scope, shows, $rootScope, $location) ->
    shows.getAll (err, shows) ->
        if err?
            $rootScope.error = err.msg
        else
            $scope.shows = shows

    $scope.showSeasons = (index) ->
        shows.setShow $scope.shows[index]
        $location.path 'shows/seasons'

]

raspTv.controller 'seasonsCtrl', ['$scope', '$rootScope', '$location', 'shows', ($scope, $rootScope, $location, shows) ->
    $location.path('shows') if shows.getShow().length is 0

    $scope.show = shows.getShow()

    shows.getSeasons (err, seasons) ->
        if err?
            $rootScope.error = err.msg
        else
            $scope.seasons = seasons

    $scope.showEpisodes = (season) ->
        shows.setSeason season
        $location.path 'shows/seasons/episodes'
]

raspTv.controller 'episodesCtrl', ['$scope', '$rootScope', '$location', 'shows', 'player', ($scope, $rootScope, $location, shows, player) ->
    if shows.getSeason().length is 0
        $location.path 'shows/seasons'
    else if shows.getShow().length is 0
        $location.path 'shows'

    $scope.show = shows.getShow()
    $scope.season = shows.getSeason()

    shows.getEpisodes (err, episodes) ->
        if err?
            $rootScope.error = err.msg
        else
            $scope.episodes = episodes

    $scope.play = (episode) ->
        player.playShow $scope.show, $scope.season, episode, (err) ->
            if err?
                $rootScope.error = err.msg
            else
                shows.setEpisode episode.name
                $location.path 'play'
]