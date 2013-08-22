raspTv = angular.module 'raspTv', ['btford.socket-io']

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
]

raspTv.controller 'navCtrl', ['$scope', '$location', ($scope, $location) ->
    $scope.isActive = (page) ->
        if page is 'movies' and $location.path() is '/'
            'active'
        else if page is 'shows' and $location.path() is '/shows'
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

raspTv.controller 'movieCtrl', ['$scope', '$http', '$rootScope', '$location', ($scope, $http, $rootScope, $location) ->
    req = $http.get '/movies'
    req.success (data) ->
        $scope.movies = data.movies
    req.error (err) ->
        $rootScope.error = err.msg

    $scope.play = (index) ->
        req = $http.post '/movies/play',
            movie : $scope.movies[index]
        req.success () ->
            $rootScope.playing = $scope.movies[index]
            $rootScope.isPlaying = true
            $rootScope.isPaused = false
            $location.path 'play'
        req.error (err) ->
            $rootScope.error = err.msg
]

raspTv.controller 'playCtrl', ['$scope', '$location', 'socket', ($scope, $location, socket) ->
    $location.path('/') if not $scope.isPlaying

    $scope.pause = () ->
        socket.emit 'pause', {paused : $scope.isPaused}, (data) ->
            $scope.isPaused = true
    $scope.play = () ->
        socket.emit 'play', (data) ->
    $scope.backward = () ->
        socket.emit 'backward', (data) ->
    $scope.forward = () ->
        socket.emit 'forward', (data) ->
]

raspTv.controller 'showsCtrl', ['$scope', ($scope) ->
]