package main

import (
    "golang.org/x/net/websocket"
    "fmt"
    "log"
    "net/http"
    "sync"
    "encoding/json"
)

type ChatMsg struct {
    Pseudo string
    Msg string
}

type Response struct {
    Recognizer string
    Pseudo string
    Msg string
}

type LoadUsersConnected struct {
    Recognizer string
    AllUserConnected []string
}

type AllMsgSend struct {
    Recognizer string
    AllMsg []ChatMsg
}

type UserDeconnected struct {
    Recognizer string
    Pseudo string
}


var USERS_LIST map[*websocket.Conn]string
// var USERS_MSG map[string]string
var USERS_MSG []ChatMsg

func Echo(ws *websocket.Conn) {
    var err error
    var m = &sync.Mutex{}

    // Associer la websocket de l'user qui vient de se connecter a un nom 
    // Stocker les messages et les noms d'utilisateurs dans une liste ?
    // Si il y a des utilisateurs dans la liste envoyer tout les messages
    // Lorsque qu'un utilisateur se deconnecte le supprimer de la liste (a faire avant la boucle)

    for {
        var reply string

        fmt.Println(len(USERS_LIST))

        err = websocket.Message.Receive(ws, &reply);
        if err != nil {
        
            if len(USERS_LIST) > 0 {
                for wsuser, username := range USERS_LIST {
                    if wsuser == ws {
                        ws.Close()
                        delete(USERS_LIST, wsuser)

                        userDisconnected := &UserDeconnected{
                           "userDisconnected",
                           username,
                        }

                        for userws := range USERS_LIST {
                            err = websocket.JSON.Send(userws, userDisconnected);
                            if err != nil {
                                fmt.Println("Can't send user disconnected")
                                break
                            }
                        }
                    }
                }
            }

            fmt.Println("Can't receive")
            break
        } 

        // JSON DECODE DE LA SOCKET + VERIFICATION SI C'EST UN NOUVEL UTILISATEUR
        r := &Response{}
        err := json.Unmarshal([]byte(reply), &r)
        if err != nil {
            panic(err)
        }

        if r.Recognizer == "newuser" {

            // SI IL Y A AU MOINS 1 USER CONNECTED RECUPERER LA LISTE ET L'ENVOYER
            if len(USERS_LIST) > 0 {
                
                var allUserConnected = []string{}

                for _, user := range USERS_LIST {
                    allUserConnected = append(allUserConnected, user)
                }

                resAllUser := &LoadUsersConnected{
                    "loadUsersConnected",
                    allUserConnected,
                }

                err = websocket.JSON.Send(ws, resAllUser);
                if err != nil {
                    fmt.Println("Can't send list of users connected")
                    break
                }
             
                if len(USERS_MSG) > 0 {

                    resAllMsgSend := &AllMsgSend{
                        "AllMsgSend",
                        USERS_MSG,
                    }

                    err = websocket.JSON.Send(ws, resAllMsgSend);
                    if err != nil {
                        fmt.Println("Can't send all msg of users")
                        break
                    }

                }

            }

            // ENVOI DE LA CONNEXION DU NOUVEL UTILISATEUR COTE CLIENT
            m.Lock()
            USERS_LIST[ws] = r.Pseudo
            m.Unlock()

            resforclient := &Response{
                "newuser",
                r.Pseudo,
                "",
            }

            for userws := range USERS_LIST {
                err = websocket.JSON.Send(userws, resforclient);
                if err != nil {
                    fmt.Println("Can't send")
                    break
                } 
            }

        } 

        if r.Recognizer == "newmessage" {

            chatMsg := &ChatMsg{
                r.Pseudo,
                r.Msg,
            }

            // SAVE MSG IN USER_MSG
            m.Lock()
            USERS_MSG = append(USERS_MSG, *chatMsg)
            m.Unlock()
            
            // SEND MESSAGE TO CLIENT
            resforclient := &Response{
                "newmessage",
                r.Pseudo,
                r.Msg,
            }

            for userws := range USERS_LIST {
                err = websocket.JSON.Send(userws, resforclient);
                if err != nil {
                    fmt.Println("Can't send")
                    break
                } 
            }

            fmt.Println(USERS_MSG)

        }


    }

}



func main() {
    USERS_LIST = make(map[*websocket.Conn]string)
    // USERS_MSG = make(map[string]string)
    // USERS_MSG = make([]ChatMsg)
    http.Handle("/", websocket.Handler(Echo))

    if err := http.ListenAndServe(":8000", nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

