package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var channel = ""

type Response struct {
	w      http.ResponseWriter
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (r *Response) SendJson() (int, error) {
	resp, _ := json.Marshal(r)
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(r.Status)
	return r.w.Write(resp)
}

func newFunc(s string) {
	go startDownload(s)
}

func handler(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}

	defer res.SendJson()

	if err := r.ParseForm(); err != nil {
		res.Status = 400
		res.Data = "url: /data?channel=<ROOM_NAME>"
		return
	}

	if r.URL.Query().Get("channel") == "$tart" {

		res.Data = "dfg"
		res.Status = 200
		response := fmt.Sprintf("Response for channel=%s", channel)
		fmt.Fprint(w, response)

		newFunc(channel)
		return
	} else {
		channel = r.URL.Query().Get("channel")
	}

	if len(channel) == 0 {
		res.Status = 400
		res.Data = "url: /data?channel=<ROOM_NAME>"
		return
	}

	fmt.Println(channel)

	res.Data = "ok"
	fmt.Print("After res.data")

	response := fmt.Sprintf("Response for channel=%s", channel)

	fmt.Fprint(w, response)

	// duration := 60 * time.Second
	// time.Sleep(duration)
	fmt.Print("before vlc")

	// // apiURL := "https://localhost:7001/data?channel=" + channel
	// res1, err := http.Get(apiURL)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// defer res1.Body.Close()
}

// func startApi() {
// 	http.HandleFunc("/", handler)
// 	http.HandleFunc("/download", handler2)

// 	// Start the HTTP server on port 8080
// 	err := http.ListenAndServe(":7001", nil)
// 	if err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// 	fmt.Println("Server Started")

// }

func startDownload(s string) {

	fmt.Printf("download started")
	savePath := "/Users/vatsalpatel/Downloads/segments/" + s
	url := "rtmp://localhost:1935/appname/" + s

	err := os.Mkdir(savePath, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	pathName := "/Users/vatsalpatel/Downloads/segments/" + s + "/index.m3u8"
	go callAPI()
	cmd := exec.Command("ffmpeg", "-i", url, "-c:v", "libx264", "-c:a", "aac", "-f", "hls", "-hls_time", "10", "-hls_list_size", "6", "-hls_flags", "6", "-vsync", "1", pathName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}

func callAPI() {
	apiURL2 := "http://localhost:9001/ded?channel=" + channel
	fmt.Println(apiURL2)
	res2, err := http.Get(apiURL2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res2.Body.Close()
}

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":7001", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	fmt.Println("Server Started")
	// go startApi()

	// cmd := exec.Command("ffmpeg", "-i", "rtmp://localhost:1935/appname/channel5", "-c:v", "libx264", "-c:a", "aac", "-f", "hls", "-hls_time", "10", "-hls_list_size", "6", "-hls_flags", "6", "-vsync", "1", "/Users/vatsalpatel/Downloads/segments/index.m3u8")

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}
