package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/andlabs/ui"
)

var KNOWN_PORTS = map[int]string{
	27017: "mongodb [ http://www.mongodb.org/ ]",
	28017: "mongodb web admin [ http://www.mongodb.org/ ]",
	21:    "ftp",
	22:    "SSH",
	23:    "telnet",
	25:    "SMTP",
	66:    "Oracle SQL*NET?",
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
	9160:  "Cassandra [ http://cassandra.apache.org/ ]",
}

var puertos_encontrados = make(map[int]string)
var host string
var puerto_inicial, puerto_final int
var wait_group sync.WaitGroup
var progreso int
var numero_de_puertos int
var mutex = sync.Mutex{}

const TIEMPO_DE_ESPERA = 15 * 1000 * time.Millisecond //milisegundos

func main() {
	err := ui.Main(func() {
		input := ui.NewEntry()
		button := ui.NewButton("Greet")
		greeting := ui.NewLabel("")
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("Enter your name:"), false)
		box.Append(input, false)
		box.Append(button, false)
		box.Append(greeting, false)
		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			greeting.SetText("Hello, " + input.Text() + "!")
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}

	pedir_info_conexion()
	obtener_puertos_abiertos()
	fmt.Println(puertos_encontrados)
}

func obtener_puertos_abiertos() {
	numero_de_puertos = puerto_final - puerto_inicial + 1
	// puertos_por_hilo := int(numero_de_puertos / 100)
	wait_group.Add(numero_de_puertos)
	// if puertos_por_hilo > 1 {
	// 	residuo := numero_de_puertos % 100
	// 	tmp_puerto_inicial := puerto_inicial
	// 	for i := 1; i <= 100; i++ {
	// 		puerto_extra := 0
	// 		if residuo != 0 {
	// 			puerto_extra = 1
	// 			residuo--
	// 		}
	// 		go revisar_puertos(tmp_puerto_inicial, tmp_puerto_inicial+puertos_por_hilo-1+puerto_extra)
	// 		tmp_puerto_inicial += puertos_por_hilo + puerto_extra
	// 	}
	// } else {
	for puerto_a_revisar := puerto_inicial; puerto_a_revisar <= puerto_final; puerto_a_revisar++ {
		go revisar_puerto(puerto_a_revisar)
	}
	// }
	wait_group.Wait()
}

func pedir_info_conexion() {
	fmt.Print("URL de servidor> ")
	fmt.Scan(&host)
	fmt.Print("Desde qué puerto (ejemplo: 80)> ")
	fmt.Scan(&puerto_inicial)
	fmt.Print("Hasta qué puerto (ejemplo 8080)> ")
	fmt.Scan(&puerto_final)
	for puerto_final < puerto_inicial {
		fmt.Println("El puerto final no puede ser menor que el puerto inicial")
		fmt.Print("Desde qué puerto (ejemplo: 80)> ")
		fmt.Scan(&puerto_inicial)
		fmt.Print("Hasta qué puerto (ejemplo 8080)> ")
		fmt.Scan(&puerto_final)
	}
	fmt.Println("E S C A N E A N D O    P U E R T O S... ")
}

func revisar_puertos(puerto_inicial, puerto_final int) {
	for puerto_a_revisar := puerto_inicial; puerto_a_revisar <= puerto_final; puerto_a_revisar++ {
		revisar_puerto(puerto_a_revisar)
	}
}

func revisar_puerto(puerto int) {
	servidor := fmt.Sprintf("%s%s%d", host, ":", puerto)
	conn, err := net.DialTimeout("tcp", servidor, TIEMPO_DE_ESPERA)
	if err != nil {
		progreso++
		fmt.Printf("\rProgeso %d/%d", progreso, numero_de_puertos)
		wait_group.Done()
		return
	}
	conn.Close()
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n1\n22\n\n\n\n")
	//conn.SetReadDeadline(time.Now().Add(TIEMPO_DE_ESPERA))
	//buff := make([]byte, 1024)
	//n, _ := conn.Read(buff)
	//fmt.Printf("Port: %d%s\n", puerto, buff[:n])
	//content := string(buff[:n])
	mutex.Lock()
	puertos_encontrados[puerto] = "" //content
	mutex.Unlock()
	progreso++
	fmt.Printf("\rProgeso %d/%d", progreso, numero_de_puertos)
	wait_group.Done()
}
