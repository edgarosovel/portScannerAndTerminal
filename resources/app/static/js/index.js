const PUERTOS_CONOCIDOS = {
	27017: "mongodb",
	28017: "mongodb web admin",
	21:    "ftp",
	22:    "SSH",
	23:    "telnet",
	25:    "SMTP",
	66:    "Oracle SQL*NET",
	69:    "tftp",
	80:    "http",
	88:    "kerberos",
	109:   "pop2",
	110:   "pop3",
	123:   "ntp",
	137:   "netbios",
	139:   "netbios",
	443:   "https",
	445:   "Samba",
	631:   "cups",
	5800:  "VNC remote desktop",
	194:   "IRC",
	118:   "SQL service?",
	150:   "SQL-net?",
	1433:  "Microsoft SQL server",
	1434:  "Microsoft SQL monitor",
	3306:  "MySQL",
	3396:  "Novell NDPS Printer Agent",
	3535:  "SMTP (alternate)",
	554:   "RTSP",
	9160:  "Cassandra",
}

// const FUNCIONES_ESPECIALES = {
//     22: index.iniciarSSH
// }


var enConexion = false

let index = {
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
            document.getElementById("input-terminal")
                .addEventListener("keyup", function(event) {
                if (event.keyCode === 13) {
                    index.getTerminalInput();
                }
            });
        })
    },
    agregarBotonPuerto: function(puerto){
        let boton = document.createElement("input");
        boton.className = "btn btn-success btn-block";
        boton.value = PUERTOS_CONOCIDOS[puerto] ? puerto + " ("+PUERTOS_CONOCIDOS[puerto]+")" : puerto;
        boton.type = "submit";
        boton.onclick = function() { index.conectarAPuerto(puerto) };
        document.getElementById("botones-puertos").appendChild(boton);
    },
    escanearPuertos: function(){
        let message = {"name": "escanearPuertos"};

        // OBTENER PARAMETROS DE INPUTS
        var servidor = document.getElementById("servidor").value;
        var desde = Number(document.getElementById("desde").value);
        var hasta = Number(document.getElementById("hasta").value);
        
        if(servidor=="" || desde=="" || hasta==""){
            asticode.notifier.error("Por favor llene todos los campos.");
            return
        }
        if(desde<1){
            asticode.notifier.error("El puerto de inicio debe ser 1 o mayor.");
            return
        }

        if(desde>hasta){
            asticode.notifier.error("El puerto de inicio es mayor que el final.");
            return
        }

        if(hasta>100000){
            asticode.notifier.error("Puerto final muy grande (Mayor de 100,000).");
            return
        }

        // TODO: regex para revisar el servidor

        message.payload = {
            servidor: servidor,
            desde: desde,
            hasta: hasta
        };

        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            if (message.name === "error") {
                asticode.notifier.error(message.name);
                return
            }
            document.getElementById("botones-puertos").innerHTML = ""

            if(message.payload.puertos.length==0){
                astilectron.showMessageBox({message:"No se encontraron puertos abiertos", title: "Puertos cerrados"})
                return
            }

            // agregar botones de puertos encontrados
            for (var puerto of message.payload.puertos){
                index.agregarBotonPuerto(puerto.puerto);
            }
        })
    },
    conectarAPuerto: function(puerto){
        let message = {"name": "conectarAPuerto"};
        message.payload = puerto;
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            if (message.name === "error") {
                asticode.notifier.error(message.name);
                return
            }
            if (message.payload.error){
                asticode.notifier.error(message.payload.respuesta);
                return
            }
            enConexion=true
            titulo = "Terminal | Puerto: "+puerto
            if(PUERTOS_CONOCIDOS[puerto]) titulo += " ("+PUERTOS_CONOCIDOS[puerto]+")";
            document.getElementById("titulo-terminal").textContent = titulo;
            document.getElementById("terminal-contenido").innerHTML = message.payload.respuesta+"<br>";
        })
    },
    iniciarSSH(){
        //PEDIR Y OBTENER USUARIO Y CONTRASEÑA
        puerto = 22
        titulo = "Terminal | Puerto: "+puerto
        if(PUERTOS_CONOCIDOS[puerto]) titulo += " ("+PUERTOS_CONOCIDOS[puerto]+")";
        document.getElementById("titulo-terminal").textContent = titulo;
        document.getElementById("terminal-contenido").innerHTML = "Ingrese usuario SSH:"+"<br>";

        let message = {"name": "iniciarSSH"};
        message.payload = {
            usuario: usuario,
            pass: pass
        };
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            if (message.name === "error") {
                asticode.notifier.error(message.name);
                return
            }
            if (message.payload.error){
                asticode.notifier.error(message.payload.respuesta);
                return
            }
            enConexion=true
            titulo = "Terminal | Puerto: "+puerto
            if(PUERTOS_CONOCIDOS[puerto]) titulo += " ("+PUERTOS_CONOCIDOS[puerto]+")";
            document.getElementById("titulo-terminal").textContent = titulo;
            document.getElementById("terminal-contenido").innerHTML = message.payload.respuesta+"<br>";
        })
    },
    getTerminalInput(){
        input = document.getElementById("input-terminal").value;
        document.getElementById("input-terminal").value = "";
        if(enConexion){
            document.getElementById("terminal-contenido").innerHTML += "💻> "+input + "<br>";
            let message = {"name": "mensajeAPuerto"};
            message.payload = input;
            asticode.loader.show();
            astilectron.sendMessage(message, function(message) {
                // Init
                asticode.loader.hide();
                if (message.name === "error") {
                    asticode.notifier.error(message.name);
                    return
                }
                if (message.payload.error){
                    enConexion = false
                    asticode.notifier.error(message.payload.respuesta);
                    document.getElementById("terminal-contenido").innerHTML = "";
                    document.getElementById("titulo-terminal").textContent = "Terminal | Sin conexión";
                    return
                }
                var terminalContenido = document.getElementById("terminal-contenido");
                terminalContenido.innerHTML += message.payload.respuesta+"<br>";
                terminalContenido.scrollTo(0,terminalContenido.scrollHeight);
            })
        }
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
                case "progreso":
                document.getElementById("botones-puertos").innerHTML =
                '<div class="progress"><div class="progress-bar progress-bar-striped active" role="progressbar" style="width:'+message.payload+'%">'+message.payload+'%</div></div>';
            }
        });
    }
};