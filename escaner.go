package main

import (
	"fmt"
	"net"
	"sync"
	"time"
	"github.com/asticode/go-astilectron-bootstrap"
)

var puertos_encontrados []int
var host string
var numero_de_puertos int
var progreso int
var wait_group sync.WaitGroup
var mutex = sync.Mutex{}
var conexion net.Conn

const TIEMPO_DE_ESPERA = 15 * 1000 * time.Millisecond //milisegundos

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

func MensajeAPuerto(input string) (bool,string) {
	if (conexion!=nil){
		fmt.Fprintf(conexion, input+"\n")
		conexion.SetReadDeadline(time.Now().Add(TIEMPO_DE_ESPERA))
		buff := make([]byte, 1024)
		n, err := conexion.Read(buff)
		if (err != nil){
			return true, "Se perdió la conexión al puerto / servidor"
		}
		response := string(buff[:n])
		return false, response
	}
	return true, "Se perdió la conexión al puerto / servidor"
}

func conectarAPuerto(puerto int) (bool,string) {
	servidor := fmt.Sprintf("%s%s%d", host, ":", puerto)
	var err error
	conexion, err = net.DialTimeout("tcp", servidor, TIEMPO_DE_ESPERA)
	if(err!=nil){
		return true,"No se pudo realizar la conexión al puerto"
		conexion = nil
	}
	return false,"Conectado"
}

func cerrarConexion()(bool,string){
	conexion.Close()
	conexion = nil
	return false,"Desconectado"
}
