package slack

import (
    "log"
    "golang.org/x/net/websocket"
    "net/http"
    "net/url"
    "encoding/json"
    "io"
    "os"
)

const BOT_NAME = "imghr"

type RtmStart struct {
    Url string
}

type AbstractRestful struct {
    Ok bool
}

func ConnectSocket(token string) *websocket.Conn {
    v := url.Values{}
    v.Set("token", token)
    response, _ := http.PostForm("https://slack.com/api/rtm.start", v)
    dec := json.NewDecoder(response.Body)
    var data RtmStart
    dec.Decode(&data)
    log.Print("start connect to ", data.Url)
    ws, err := websocket.Dial(data.Url, "", "http://localhost")
    if err != nil {
        log.Fatal(err)
    }

    return ws
}

func StartReading(conn *websocket.Conn) (<-chan []byte, <-chan bool) {
    log.Print("start reading...")
    var msg []byte
    sendChan := make(chan []byte)
    breakChan := make(chan bool)

    go func() {
        for {
            var tmp = make([]byte, 512)
            n, err := conn.Read(tmp)
            if err == io.EOF {
                breakChan <- true
                break
            }
            if err != nil {
                log.Fatal(err)
            }
            if msg != nil {
                msg = append(msg, tmp[:n]...)
            } else {
                msg = tmp[:n]
            }
            if n != 512 {
                sendChan <- msg
                msg = nil
            }
        }
    }()

    return sendChan, breakChan
}

func PostMessage(channel string, username string, text string) <-chan bool {
    token := os.Getenv("SLACK_TOKEN")
    resultChan := make(chan bool)
    v := url.Values{}
    v.Set("token", token)
    v.Set("channel", channel)
    v.Set("text", text)
    v.Set("username", username)
    go func() {
        response, _ := http.PostForm("https://slack.com/api/chat.postMessage", v)
        dec := json.NewDecoder(response.Body)
        var data AbstractRestful
        dec.Decode(&data)
        resultChan <- data.Ok
    }()

    return resultChan
}
