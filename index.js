(function() {

    var sock = null;
    var wsuri = "ws://localhost:8000/ws";
    var pseudoUser = "";

    window.onload = function() {

        document.getElementById('btn').addEventListener("click", function() {
            var pseudo = document.getElementById('pseudo').value;
            
            if (pseudo.length > 2) {
                document.getElementById('hide-block').style.display = "none";

                sock = new WebSocket(wsuri);

                // Cet événement se produit lorsque la connexion socket est établie.
                sock.onopen = function() {
                    console.log("connected to " + wsuri);
                    var newUserPseudo = { "recognizer": "newuser", "pseudo": pseudo, "msg": "" };
                    sock.send(JSON.stringify(newUserPseudo));
                }

                // Cet événement se produit lorsque la connexion est fermée.
                sock.onclose = function(e) {
                    console.log("connection closed (" + e.code + ")");
                }

                // Retour de la Reponse Serveur 
                sock.onmessage = function(e) {
                    console.log("Reponse API -> msg : " + e.data);
                    var resFromApi = JSON.parse(e.data)

                    if (resFromApi["Recognizer"] == "loadUsersConnected") {
                        var data = JSON.parse(e.data)

                        for (pseudo in data["AllUserConnected"])
                        {
                            var userPseudo = data["AllUserConnected"][pseudo];
                            var iDiv = document.createElement('div');
                            iDiv.id = userPseudo;
                            iDiv.className = 'user';
                            iDiv.innerHTML = userPseudo;
                            document.getElementById("user-connected").appendChild(iDiv);
                        }
                    }

                    if (resFromApi["Recognizer"] == "AllMsgSend") {
                        var data = JSON.parse(e.data)
                        for (key in data["AllMsg"])
                        {
                            var usrPseudo = data["AllMsg"][key]["Pseudo"];
                            var usrMessage = data["AllMsg"][key]["Msg"];

                            var newPMsg = document.createElement('p');
                            newPMsg.innerHTML = usrPseudo + " : " + usrMessage
                            document.getElementById("block-response").appendChild(newPMsg);
                        }
                    }

                    if (resFromApi["Recognizer"] == "newuser") {
                        pseudoUser = resFromApi["Pseudo"];
                    
                        var iDiv = document.createElement('div');
                        iDiv.id = resFromApi["Pseudo"];
                        iDiv.className = 'user';
                        iDiv.innerHTML = resFromApi["Pseudo"];
                        document.getElementById("user-connected").appendChild(iDiv);

                    } 

                    if (resFromApi["Recognizer"] == "newmessage") {
                        var newPMsg = document.createElement('p');
                        newPMsg.innerHTML = resFromApi["Pseudo"] + " : " + resFromApi["Msg"]
                        document.getElementById("block-response").appendChild(newPMsg);
                    }

                    if (resFromApi["Recognizer"] == "userDisconnected") {
                        var data = JSON.parse(e.data);
                        var pseudoDeleted = data["Pseudo"];
                        var divToDel = document.getElementById(pseudoDeleted);
                        divToDel.remove();
                    }
                    
                }

                document.getElementById('btn-submit').addEventListener("click", function() {
                    var msg = document.getElementById('message').value;
                    var newMsg = { "recognizer": "newmessage", "pseudo": pseudoUser, "msg": msg };
                    sock.send(JSON.stringify(newMsg));
                });

            } 
            else 
            {
                console.log('enter valid pseudo');
            }


        });

    };

    

})();