window.onload = function() {
    let statusIcon = this.document.getElementById("status-icon");
    let socket = new WebSocket("wss://" + document.location.host + "/ws");
    let keysDown = {};

    socket.onopen = function() {
        $('#status-icon').attr('data-icon', 'link')
            .addClass('text-success')
            .removeClass('text-warning');
    
        _firstStart();

        socket.send("Tadaa!");
    };

    socket.onclose = function() {
        $('#status-icon').attr('data-icon', 'unlink')
            .addClass('text-warning')
            .removeClass('text-success');
    };

    socket.onmessage = function(msg) {
        $('#last-ws-box').html('Last response: "'+msg.data+'"');

        let split = msg.data.split("|");

        if (split.length === 2) {
            let event = split[0];
            let data = split[1];

            switch (event) {
                case "player_update":
                    let gameObject = JSON.parse(data);

                    player.x = gameObject.x;
                    player.y = gameObject.y;
                    player.w = gameObject.w;
                    player.h = gameObject.h;

                    console.dir(gameObject);
                    break;
                default:
                    console.error("Can't handle event: "+event);
            }
        }
    };

    const K_UP = 38;
    const K_DOWN = 40;
    const K_LEFT = 37;
    const K_RIGHT = 39;
    const K_W = 87;
    const K_S = 83;
    const K_A = 65;
    const K_D = 68;

    const X_SPEED = 40;
    const Y_SPEED = 30;

    const canvas = this.document.getElementById("game-canvas");
    let ctx;
    let lastTime;
    const cw = canvas.width;
    const ch = canvas.height;
    let player;
    let objects = [];

    let paused = false;

    function _firstStart() {
        
        ctx = canvas.getContext('2d');
        player = {
            x: 30,
            y: 30,
            w: 50,
            h: 50
        }

        objects.push(player);

        lastTime = (new Date()).getTime();

        _bindKeys();
        _gameLoop();
    }

    function _bindKeys() {
        $(document).keydown(function(e) {
            var key = e.keyCode;

            keysDown[key] = true;
        });

        $(document).keyup(function(e) {
            var key = e.keyCode;
            
            keysDown[key] = false;
        })
    }

    function _gameLoop() {
        window.requestAnimationFrame(_gameLoop);

        if (paused) {
            return;
        }

        var currentTime = (new Date()).getTime();
        var delta = (currentTime - lastTime)/1000;
        lastTime = currentTime;
        
        _move(delta);
        _draw(delta);
        _sync(delta);

    }

    function _move(delta) {
        let yAxis = 0;
        if (keysDown[K_W] || keysDown[K_UP]) {
            yAxis += 1;
        }
        if (keysDown[K_S] || keysDown[K_DOWN]) {
            yAxis -= 1;
        }

        let xAxis = 0;
        if (keysDown[K_A] || keysDown[K_LEFT]) {
            xAxis -= 1;
        }
        if (keysDown[K_D] || keysDown[K_RIGHT]) {
            xAxis += 1;
        }

        player.x += xAxis * X_SPEED * delta;
        player.y += yAxis * Y_SPEED * delta;
    }

    function _draw(delta) {
        ctx.clearRect(0, 0, cw, ch)

        ctx.fillStyle = "rgb(200, 0, 0)";
        objects.forEach(obj => {
            ctx.fillRect(obj.x, obj.y, obj.w, obj.h)
        });
    }

    function _sync(delta) {
        socket.send("player_update|"+JSON.stringify(player));
    }
};
