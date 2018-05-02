package main

import (
	"assistant_game_server/client"
	"assistant_game_server/config"
	"assistant_game_server/db"
	"assistant_game_server/gameconf"
	"assistant_game_server/gameredis"
	"assistant_game_server/newspush"
	"assistant_game_server/task"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func initLogger() {
	logDir := config.Get().Log.LogDir
	os.Mkdir(logDir, os.ModeDir)
	logDirSet := false
	for _, arg := range os.Args {
		if strings.Contains(arg, "-log_dir=") {
			logDirSet = true
			break
		}
	}
	if !logDirSet {
		os.Args = append(os.Args, "-log_dir="+logDir)
	}
	flag.Parse()
}

func handleClient(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Upgrade(w, r, nil, 65536, 65536)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	client := &client.GameClient{}
	go client.HandleConnection(c)
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/meatfloss")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	rand.Seed(time.Now().UTC().UnixNano())
	var err error
	if err = config.LoadConfig(); err != nil {
		panic(fmt.Sprintf("LoadConfig failed, %s\n", err))
	}

	initLogger()
	glog.Info("starting...")

	err = db.Initialize(config.Get().MySQLServer[0].Host,
		config.Get().MySQLServer[0].Port,
		config.Get().MySQLServer[0].User,
		config.Get().MySQLServer[0].Password)
	if err != nil {
		glog.Fatalf("db.Initialize failed, %s", err)
	}

	err = gameconf.LoadFromDatabase()
	if err != nil {
		glog.Fatalf("gameconf.LoadFromDatabase failed, error: %s", err)

	}

	gameredis.Initialize()

	http.HandleFunc("/meatfloss", handleClient)
	http.HandleFunc("/", home)
	go newspush.RunHTTPServer()
	task.StartTaskManager()
	//	guest.Start()
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
