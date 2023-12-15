package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func dedHandler(w http.ResponseWriter, r *http.Request) {
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

	channel := r.URL.Query().Get("channel")
	fmt.Println(channel)

	if len(channel) == 0 {
		res.Status = 400
		res.Data = "url: /data?channel=<ROOM_NAME>"
		return
	}

	// msg, err := configure.RoomKeys.GetKey(room)
	// if err != nil {
	// 	msg = err.Error()
	// 	res.Status = 400
	// }

	response := fmt.Sprintf("Response for channel=%s", channel)

	fmt.Fprint(w, response)

	go startStream(channel)

}

func startStream(s string) {
	segDir := "/Users/vatsalpatel/downloads/segments/" + s
	const port = 9002
	fmt.Println(segDir)

	http.Handle("/", addHeaders(http.FileServer(http.Dir(segDir))))
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", segDir, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func main() {
	http.HandleFunc("/ded", dedHandler)

	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
