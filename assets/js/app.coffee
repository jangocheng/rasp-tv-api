raspTv = angular.module 'raspTv', ['btford.socket-io']

raspTv.config ['$routeProvider', ($routeProvider) ->
    $routeProvider.when '/',
        templateUrl : '/templates/movies.html'
        controller : 'movieCtrl'
]

raspTv.controller 'navCtrl', ['$scope', '$location', ($scope, $location) ->
    $scope.isActive = (page) ->
        if page is 'movies' and $location.path() is '/'
            'active'
        else if page is 'shows' and $location.path() is '/shows'
            'active'
        else
            ''
]

raspTv.controller 'errorCtrl', ['$scope', ($scope) ->
    $scope.close = () ->
        $scope.error = null
]

raspTv.controller 'movieCtrl', ['$scope', '$http', '$rootScope', ($scope, $http, $rootScope) ->
    req = $http.get '/movies'
    req.success (data) ->
        $scope.movies = data.movies
    req.error (err) ->
        $rootScope.error = err.msg

    $scope.play = (index) ->
        req = $http.post '/movies/play',
            movie : $scope.movies[index]
        req.success () ->
            console.log 'movie play'
        req.error (err) ->
            $rootScope.error = err.msg
]