window.onload = function() {

    var socket = new WebSocket("wss://" + document.location.host + "/ws");
    var msg = document.getElementById("msg");

    document.getElementById("form").onsubmit = function () {
        if (!socket) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        socket.send(msg.value);
        msg.value = "";
        return false;
    };
};