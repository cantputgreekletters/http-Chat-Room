package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const message_buffer_count int = 10

type Tail struct {
	data        [message_buffer_count]string
	current_idx int
	count       int
}

func (t *Tail) push(message string) {
	t.current_idx = (t.current_idx + 1) % message_buffer_count
	t.data[t.current_idx] = message

	if t.count < message_buffer_count {
		t.count++
	}
}

func (t *Tail) get() [message_buffer_count]string {
	var data [message_buffer_count]string

	start := (t.current_idx - t.count + 1 + message_buffer_count) % message_buffer_count

	for i := 0; i < t.count; i++ {
		data[i] = t.data[(start+i)%message_buffer_count]
	}

	return data
}

type Message struct {
	Type string                       `json:"type"`
	Data [message_buffer_count]string `json:"data"`
}

var Messages Tail

func update_api_handling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	msg := Message{
		Type: "chat",
		Data: Messages.get(),
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		println("Failed to send update")
		return
	}
	w.Write(jsonData)
}

func message_api_handling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Printf("Error = %s\n", err.Error())
		return
	}

	message := string(body)
	if message == "" {
		return
	}
	fmt.Println("Received:", message)
	Messages.push(message)
	w.Write([]byte(message))
}

//func update_api_handling

func Setup(URL string) {
	url := URL + ":8080"
	Messages.current_idx = -1
	fs := http.FileServer(http.Dir("./front"))
	http.Handle("/", fs)

	http.HandleFunc("/api/message", message_api_handling)
	http.HandleFunc("/api/update", update_api_handling)

	fmt.Printf("Serving on url %s\n", url)
	var err error = http.ListenAndServe(url, nil)
	if err != nil {
		fmt.Printf("Error = %s\n", err.Error())
		return
	}
}
