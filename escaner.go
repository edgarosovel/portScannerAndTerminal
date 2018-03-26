package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"regexp"
	"io"
	"time"
	"bytes"
	"github.com/asticode/go-astilectron-bootstrap"
	"golang.org/x/crypto/ssh"
	"github.com/mitchellh/go-homedir"
	// "strconv"
)

var puertos_encontrados []int
var host string
var numero_de_puertos int
var progreso int
var wait_group sync.WaitGroup
var mutex = sync.Mutex{}

var conexion net.Conn
var sesionSSH *ssh.Session

const TIEMPO_DE_ESPERA = 7 * 1000 * time.Millisecond //milisegundos

func obtener_puertos_abiertos(_host string, puerto_inicial int, puerto_final int) []int {
	host = _host
	progreso=0
	puertos_encontrados = []int{}
	numero_de_puertos = puerto_final - puerto_inicial + 1
	wait_group.Add(numero_de_puertos)
	for puerto_a_revisar := puerto_inicial; puerto_a_revisar <= puerto_final; puerto_a_revisar++ {
		go revisar_puerto(puerto_a_revisar)
	}
	wait_group.Wait()
	return puertos_encontrados
}

type MensajeProgreso struct{
	Name string `json:"name"`
	Payload int `json:"payload"`
}

func revisar_puerto(puerto int) {
	servidor := fmt.Sprintf("%s%s%d", host, ":", puerto)
	conn, err := net.DialTimeout("tcp", servidor, TIEMPO_DE_ESPERA)
	if err != nil {
		progreso++
		if err := bootstrap.SendMessage(w, "progreso",  progreso * 100 / numero_de_puertos); err != nil {	}
		wait_group.Done()
		return
	}
	conn.Close()
	mutex.Lock()
	puertos_encontrados = append(puertos_encontrados, puerto)
	mutex.Unlock()
	progreso++
	if err := bootstrap.SendMessage(w, "progreso",  progreso * 100 / numero_de_puertos); err != nil {}
	wait_group.Done()
}

func mensajeAPuerto(input string) (bool,string) {
	if conexion!=nil{
		fmt.Fprintf(conexion, input+"\n")
		conexion.SetReadDeadline(time.Now().Add(TIEMPO_DE_ESPERA))
		var buff bytes.Buffer
		io.Copy(&buff, conexion) //leer datos
		if buff.Len()==0 {
			return true, "Se perdió la conexión al puerto / servidor"
		}
		response := buff.String()
		// revisar si es archivo para descargarlo
		archivoRegex := regexp.MustCompile(`(?i)\w+\.(jpg|png|gif|bmp|jpeg|pdf|zip)`)
		if esArchivo := archivoRegex.MatchString(input); esArchivo {
			res, err := guardarArchivo(buff, archivoRegex.FindString(input));
			if err != nil{
				return true, "Error al descargar el archivo"
			}
			return false, res
		}
		return false, response
	}
	return true, "Se perdió la conexión al puerto / servidor"
}

func conectarAPuerto(puerto int) (bool,string) {
	servidor := fmt.Sprintf("%s%s%d", host, ":", puerto)
	var err error
	conexion, err = net.DialTimeout("tcp", servidor, TIEMPO_DE_ESPERA)
	if err!=nil {
		conexion = nil
		return true,"No se pudo realizar la conexión al puerto"
	}
	return false,"Conectado"
}

func cerrarConexion()(bool,string){
	conexion.Close()
	sesionSSH.Close()
	return false,"Desconectado"
}


// SSH
func iniciarSSH(usuario string, pass string) (bool, string){
	var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: usuario,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	conexionSSH, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		conexionSSH = nil
		return true,"No se pudo realizar la conexión al puerto"
	}
	sesionSSH, err = conexionSSH.NewSession()
	if err != nil {
		sesionSSH.Close()
		conexionSSH.Close()
		return true,"No se pudo realizar la conexión al puerto"
	}

	return false,"Conectado"
}

func mensajeSSH(input string)(bool, string){
	if sesionSSH!=nil{
		var response bytes.Buffer
		sesionSSH.Stdout = &response
		if err := sesionSSH.Run(input); err != nil {
			return true, "Se perdió la conexión al puerto / servidor"
		}
		return false, response.String()
	}
	return true, "Se perdió la conexión al puerto / servidor"
}

func guardarArchivo(buff bytes.Buffer, fileName string)(string, error){
	homePath, _ := homedir.Dir()
	homePath, _ = homedir.Expand(homePath)
	filePath := homePath+"/"+fileName
	file, err := os.Create(filePath)
    if err != nil {
        return "Error", err
	}
	endOfHeader := bytes.Index(buff.Bytes(), []byte("\n\r\n"))
	// endOfHeader := bytes.Index(buff.Bytes(), []byte("\n"))
	// fileWithoutHeaders := bytes.SplitN(buff.Bytes(), []byte("\n"), 13)
	// var s string
	// for _, x := range fileWithoutHeaders {
	// 	fmt.Println("new line")
	// 	fmt.Println(x[:])
	// 	// s += string(x[:]) + "|"
	// }
	var fileWithoutHeaders bytes.Buffer
	fileWithoutHeaders.Write(buff.Bytes()[endOfHeader+3:])
	fileWithoutHeaders.WriteTo(file)
	file.Close()
    return "<br>Archivo guardado en: "+filePath, nil
}

// func main (){
// 	obtener_puertos_abiertos("google.com", 80, 80)
// 	conectarAPuerto(80)
// 	mensajeAPuerto("GET /logos/classicplus.png")	
// }