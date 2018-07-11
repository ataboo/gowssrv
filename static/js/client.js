var socket = new WebSocket("wss://"+document.location.host+"/ws");

socket.send("Sent from client!");