/**
 * Created by cj on 2016/12/19.
 */
var app = angular.module("chat", []);

app.directive("chat", function($location, $anchorScroll){
    return {
        link: function(scope, element, attrs){
            $location.hash('bottom');
            scope.$watch("log", function(){
                $anchorScroll();
            }, true);
        }
    }
});

app.controller("MainCtl", function ($scope, $interval) {
    $scope.login =false;
    $scope.message = "";

    if (!window["WebSocket"]) {
        $scope.log.push("Your browser does not support WebSockets.");
        return;
    }
    var conn = new WebSocket("ws://localhost:4000/ws");
    conn.onclose = function (e) {
        $scope.$apply(function () {
            $scope.log.push("Connection closed.");
        })
    };

    conn.onmessage = function (e) {
        $scope.$apply(function () {
            data = JSON.parse(e.data);
            console.log(data);
            if(data.type == "init") {
                $scope.node = data.node_addr;
                $scope.timestamp = data.node_time;
                var timer = $interval(function(){
                    $scope.timestamp += 1;
                },1000);
            }
        })
    };

    conn.onopen = function (e) {
        console.log("Connected");
    };

    $scope.send = function () {
        if (!conn) {
            return;
        }

        if (!$scope.message) {
            return;
        }

        conn.send(nick + ": " + $scope.message);
        $scope.message = "";
    }

    $scope.connect = function () {
        $scope.login = true;
        conn.send("LOGIN")
    }
});
