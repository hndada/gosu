package gosu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gorilla/websocket"
	"github.com/hndada/gosu/scene"
)

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

var musicData = []byte(`[
  {
    "musicName": "triangle!",
    "chartName": "Easy",
	"chartFS": "C:/Users/hndada/Documents/GitHub/gosu/music/cYsmix - triangles",
	"chartFilename": "cYsmix - triangles (MuangMuangE) [Easy].osu"
  }
]
`)

func openWebServer(g *Game) {
	// Serve /selects page (HTML)
	http.HandleFunc("/selects", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../scene/selects/static/index.html")
	})

	// Serve song data JSON
	http.HandleFunc("/music.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// data, err := os.ReadFile("music.json")
		// if err != nil {
		// 	fmt.Printf("Failed to read music list: %v\n", err)
		// 	http.Error(w, "Failed to read song list", 500)
		// 	return
		// }
		// w.Write(data)
		w.Write(musicData)
	})

	// Redirect / to /songs
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/selects", http.StatusFound)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(g, w, r)
	})

	// Notify game that server has started
	// Should I use a channel or context instead?

	go openBrowser("http://127.0.0.1:8080/")
	fmt.Println("Server started at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func handleWS(g *Game, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("WebSocket connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)

		// Unmarshal JSON into struct
		var argsData scene.PlayArgsData
		if err := json.Unmarshal(msg, &argsData); err != nil {
			fmt.Println("JSON error:", err)
			continue
		}

		// Send to game update loop
		go func() {
			g.events <- argsData.ToPlayArgs()
		}()
	}
}
