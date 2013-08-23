raspTv = angular.module 'raspTv', ['btford.socket-io', 'LocalStorageModule']

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

raspTv.run ['$rootScope', 'localStorageService', ($rootScope, localStorageService) ->
    $rootScope.playing = localStorageService.get 'playing'
    $rootScope.isPlaying = if localStorageService.get('isPlaying') is 'true' then true else false
    $rootScope.isPaused = if localStorageService.get('isPaused') is 'true' then true else false
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

raspTv.controller 'movieCtrl', ['$scope', '$http', '$rootScope', '$location', 'localStorageService', ($scope, $http, $rootScope, $location, localStorageService) ->
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
            localStorageService.set 'playing', $rootScope.playing
            $rootScope.isPlaying = true
            localStorageService.set 'isPlaying', $rootScope.isPlaying
            $rootScope.isPaused = false
            localStorageService.set 'isPaused', $rootScope.isPaused
            $location.path 'play'
        req.error (err) ->
            $rootScope.error = err.msg
]

raspTv.controller 'playCtrl', ['$scope', '$location', 'socket', '$rootScope', 'localStorageService', ($scope, $location, socket, $rootScope, localStorageService) ->
    $location.path('/') if not $scope.isPlaying

    $scope.toggle = () ->
        socket.emit 'toggle'
        $scope.isPaused = not $scope.isPaused
        localStorageService.set 'isPaused', $scope.isPaused
    $scope.backward = () ->
        socket.emit 'backward'
    $scope.forward = () ->
        socket.emit 'forward'
    $scope.stop = () ->
        socket.emit 'stop'
        $rootScope.isPlaying = false
        localStorageService.clearAll()
        $location.path '/'
]

raspTv.controller 'showsCtrl', ['$scope', ($scope) ->
]