package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// upgrader se encarga de actualizar la peticion http a websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // aqui validamos que la conexion sea valida
	},
}

type Mensaje struct {
	Usuario string `json:"usuario"`
	Cuerpo  string `json:"cuerpo"`
}

func (m *Mensaje) String() string {
	return fmt.Sprintf("<Usuario: %s |Mensaje %s>", m.Usuario, m.Cuerpo)
}

// Hub se encarga de llevar el reigstro de clientes, ademas de los clientes registrados
type Hub struct {
	Clientes   map[*Cliente]bool
	register   chan *Cliente
	unregister chan *Cliente
	Broadcast  chan Mensaje
	mutex      sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		Clientes:   make(map[*Cliente]bool),
		register:   make(chan *Cliente),
		unregister: make(chan *Cliente),
		Broadcast:  make(chan Mensaje, 50),
	}
}

// run mantienen corriendo en la escucha de los canales
func (h *Hub) run() {
	for {
		select {
		case cliente := <-h.register:
			h.mutex.Lock()
			h.Clientes[cliente] = true // guardamos la direccion del cliente en memoria
			h.mutex.Unlock()
			// mensamos un mensaje a los demas que se conecto
			message := fmt.Sprintf("Se ha conectado %s", cliente.Username)
			fmt.Println(message)
			h.Broadcast <- Mensaje{
				Usuario: cliente.Username,
				Cuerpo:  message,
			}

		case message := <-h.Broadcast:
			for cliente := range h.Clientes {
				select {
				case cliente.Send <- message: // corroboramos si el cliente tiene espacio para el mensaje  y so lo pudo enviar
				default: // de no ser el caso va por el siguiente caso
					// se tiene que liberar la conexion
					//FALLA: Canal lleno o bloqueado
					close(cliente.Send)
					h.mutex.Lock()
					delete(h.Clientes, cliente) // quitamos del registro de cliente
					h.mutex.Unlock()
				}
			}
		case cliente := <-h.unregister:
			_, exists := h.Clientes[cliente]
			if exists {
				fmt.Println("Cliente " + cliente.Username + " removido")
				h.mutex.Lock()
				delete(h.Clientes, cliente)
				h.mutex.Unlock()
			}
		}

	}
}

func handleConnection(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// abrimos la conexion socket hecha a la ruta /ws

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error al atualizar la peticion a websocket", err)
	}

	cliente := &Cliente{
		Username: r.URL.Query().Get("username"),
		conn:     conn,
		Send:     make(chan Mensaje, 255),
	}

	hub.register <- cliente // guardamos el cliente en la tabla de direcciones

	// quedamos a la escucha de los mensajes del cliente
	go cliente.ReadMessages(hub)
	go cliente.WriteMessages()
}
