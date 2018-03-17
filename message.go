package main

import (
	"encoding/json"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

type RespuestaPuerto struct {
	Error      bool                 `json:"error"`
	Respuesta  string 				`json:"respuesta"`
}

type Escaneo struct {
	Puertos []Puerto `json:"puertos"`
}

// PayloadDir represents a dir payload
type Puerto struct {
	Nombre string `json:"nombre"`
	Puerto int `json:"puerto"`
}

type ParametrosEscaneo struct {
	Servidor string
	Desde int
	Hasta int
}

type mensajeAPuerto struct {
	Mensaje string
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "escanearPuertos":
		var parametrosEscaneo ParametrosEscaneo
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &parametrosEscaneo); err != nil {
				payload = err.Error()
				return
			}
		}
		escaneo := Escaneo{
			Puertos: []Puerto{},
		}

		puertosEncontrados := obtener_puertos_abiertos(parametrosEscaneo.Servidor, parametrosEscaneo.Desde, parametrosEscaneo.Hasta);

		for _, puerto := range puertosEncontrados{
			escaneo.Puertos = append(escaneo.Puertos, Puerto{
				Puerto: puerto,
			})
		}

		return escaneo, err;
	case "conectarAPuerto":
		var puerto int
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &puerto); err != nil {
				payload = err.Error()
				return
			}
		}
		err, respuesta := conectarAPuerto(puerto)
		return RespuestaPuerto{err,respuesta}, nil
	case "mensajeAPuerto":
		var comando string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &comando); err != nil {
				payload = err.Error()
				return
			}
		}
		err, respuesta := MensajeAPuerto(comando)
		return RespuestaPuerto{err,respuesta}, nil
	}
	return
}