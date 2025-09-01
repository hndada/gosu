package scene

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gorilla/websocket"
)

func OpenBrowser(url string) {
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
    "musicName": "Song A",
    "chartName": "Easy",
    "filePath": "/music/songA.mp3",
    "previewStart": 30,
    "previewDur": 10
  },
  {
    "musicName": "Song B",
    "chartName": "Medium",
    "filePath": "/music/songB.mp3",
    "previewStart": 45,
    "previewDur": 12
  },
  {
    "musicName": "Song C",
    "chartName": "Hard",
    "filePath": "/music/songC.mp3",
    "previewStart": 10,
    "previewDur": 15
  },
  {
    "musicName": "Song D",
    "chartName": "Expert",
    "filePath": "/music/songD.mp3",
    "previewStart": 20,
    "previewDur": 8
  }
]
`)

func OpenWebServer() {
	// Serve /selects page (HTML)
	http.HandleFunc("/selects", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "scene/selects/static/index.html")
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

	http.HandleFunc("/ws", handleWS)
	fmt.Println("WebSocket server on ws://localhost:8080/ws")

	go OpenBrowser("http://localhost:8080/")
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func handleWS(w http.ResponseWriter, r *http.Request) {
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
	}
}
