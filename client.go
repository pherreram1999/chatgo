package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Cliente gestiona la conexion con el cliente
type Cliente struct {
	Username string
	conn     *websocket.Conn
	Send     chan Mensaje
}

// ReadMessages lee los mensajes enviados por el cliente (navegador web) al servidor (hub)
func (c *Cliente) ReadMessages(hub *Hub) {
	defer func() { // al salir de la rutina cerramos y informamos el hub
		hub.register <- c // indicamos al hub que el cliente se desconecto
		c.conn.Close()    // una vez terminado de leer los mesanjes, cerramos conexion
	}()

	for { // este loop va estar siempre disponible hasta que recibe un mensaje del cliente
		var message Mensaje
		err := c.conn.ReadJSON(&message)
		if err != nil {
			// hubo una desconexion, salimos de la rutina
			log.Println(err)
			break // aqui se rompe el ciclo, cierra la conexion al salir de la funcion
		}

		// renviamos el mensaje al resto de los clientes con un broadcast
		hub.Broadcast <- message
	}
}

// WriteMessages envia los mensaje atraves de hub broacast a los clientes en la tabla de direcciones
func (c *Cliente) WriteMessages() {
	defer c.conn.Close() // cerramos la conexion al salir del goroutine
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// se ha cerrado el canal
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				// error el escribir
				return
			}
		}
	}
}
