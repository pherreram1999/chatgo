<script lang="ts" setup>

import {ref} from "vue";


enum Steps {
  Login,
  Chat
}

interface MessageResponse {
  usuario: string,
  cuerpo: string,
  color?: string
}

const host = 'ws://192.168.8.3:8081/ws'
const username = ref<string>()
const message = ref<string>()
const selectedColor = ref<string>('blue')

const colors = ['blue', 'green', 'red', 'purple', 'pink', 'orange']

const step = ref<Steps>(Steps.Login)

let connection: WebSocket

const containerRef = ref()

const newNode = (tag: string) => document.createElement(tag)



const insertMessageNode = (msg: string) => {
  if(!containerRef.value) return
  const messageBody: MessageResponse = JSON.parse(msg)

  const divParent = newNode('div')

  // Determinar si es mensaje del usuario actual, mensaje de sistema o mensaje externo
  const isOwnMessage = messageBody.usuario === username.value
  const isSystemMessage = messageBody.cuerpo.includes('se ha conectado') ||
      messageBody.cuerpo.includes('se ha unido') ||
      messageBody.cuerpo.includes('Se ha conectado') ||
      messageBody.cuerpo.includes('Se ha unido')

  if (isSystemMessage) {
    // Mensaje de sistema centrado con estilo de texto simple
    divParent.className = 'px-4 py-1 text-center'
    const systemText = newNode('span')
    systemText.innerText = messageBody.cuerpo
    systemText.className = 'text-gray-500 dark:text-gray-400 text-sm italic font-light'
    divParent.appendChild(systemText)
  } else if (isOwnMessage) {
    // Mensaje propio alineado a la derecha
    divParent.className = 'px-4 py-2 flex justify-end'
    const messageContainer = newNode('div')
    messageContainer.className = 'max-w-xs lg:max-w-md'

    const userLbl = newNode('span')
    userLbl.innerText = messageBody.usuario
    userLbl.className = `block font-semibold text-right text-sm text-${selectedColor.value}-600 dark:text-${selectedColor.value}-400 mb-1`

    const bodyText = newNode('p')
    bodyText.innerText = messageBody.cuerpo
    bodyText.className = `bg-${selectedColor.value}-500 text-white px-4 py-2 rounded-tl-2xl rounded-bl-2xl rounded-br-2xl`

    messageContainer.appendChild(userLbl)
    messageContainer.appendChild(bodyText)
    divParent.appendChild(messageContainer)
  } else {
    // Mensaje externo alineado a la izquierda
    divParent.className = 'px-4 py-2 flex justify-start'
    const messageContainer = newNode('div')
    messageContainer.className = 'max-w-xs lg:max-w-md'

    // Usar el color del mensaje si está disponible, sino usar gris por defecto
    const messageColor = messageBody.color || 'gray'

    const userLbl = newNode('span')
    userLbl.innerText = messageBody.usuario
    userLbl.className = messageColor === 'gray'
        ? 'block font-semibold text-left text-sm text-gray-600 dark:text-gray-400 mb-1'
        : `block font-semibold text-left text-sm text-${messageColor}-600 dark:text-${messageColor}-400 mb-1`

    const bodyText = newNode('p')
    bodyText.innerText = messageBody.cuerpo
    bodyText.className = messageColor === 'gray'
        ? 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-4 py-2 rounded-tr-2xl rounded-br-2xl rounded-bl-2xl'
        : `bg-${messageColor}-500 text-white px-4 py-2 rounded-tr-2xl rounded-br-2xl rounded-bl-2xl`

    messageContainer.appendChild(userLbl)
    messageContainer.appendChild(bodyText)
    divParent.appendChild(messageContainer)
  }

  containerRef.value.appendChild(divParent)
}

const handleRecieveMessage = (event: MessageEvent) => {
  insertMessageNode(event.data)
}

const sendMessage = () => {
  if(!connection) return
  connection.send(
      JSON.stringify({
        usuario: username.value,
        cuerpo: message.value,
        color: selectedColor.value
      })
  )
  message.value = ''
}

const conectar = () => {
  const uri = new URL(host)
  uri.searchParams.append('username',username.value)
  uri.searchParams.append('color',selectedColor.value)
  connection = new WebSocket(uri.toString())
  connection.onmessage = handleRecieveMessage
  step.value = Steps.Chat
}



</script>

<template>
  <div>
    <div v-if="step === Steps.Login" class="h-screen flex justify-center items-center bg-neutral-300 dark:bg-gray-900">
      <div>
        <form @submit.prevent="conectar" class="bg-white p-6 rounded shadow-2xl text-shadow-black text-shadow-lg/60 rounded-2xl border border-stone-200 w-[26rem] dark:shadow-indigo-950 dark:border-0 dark:bg-gray-800 dark:text-white dark:text-shadow-indigo-500 dark:text-shadow-lg/60">
          <h1 class="text-4xl font-semibold text-center mb-8">WebChat Express</h1>

          <div class="mb-6">
            <input type="text"
                   name="username"
                   id="username"
                   v-model="username"
                   autofocus
                   required
                   placeholder="Nombre de usuario"
                   class="bg-stone-100 w-full rounded px-4 py-2 border-stone-500 text-black focus:outline-none focus:ring-2 focus:ring-black dark:focus:ring-indigo-500 focus:border-transparent placeholder-gray-400">
          </div>

          <div class="mb-6">
            <label class="block text-sm font-medium mb-3 text-gray-700 dark:text-gray-300">
              Color de mensaje:
            </label>
            <div class="flex gap-3 justify-center flex-wrap">
              <label
                  v-for="color in colors"
                  :key="color"
                  :class="[
                  'w-8 h-8 rounded-full cursor-pointer border-2 transition-all duration-200',
                  selectedColor === color
                    ? 'border-gray-800 dark:border-white scale-110 ring-2 ring-gray-400'
                    : 'border-gray-300 dark:border-gray-600 hover:scale-105',
                  color === 'blue' ? 'bg-blue-500' : '',
                  color === 'green' ? 'bg-green-500' : '',
                  color === 'red' ? 'bg-red-500' : '',
                  color === 'purple' ? 'bg-purple-500' : '',
                  color === 'pink' ? 'bg-pink-500' : '',
                  color === 'orange' ? 'bg-orange-500' : ''
                ]">
                <input
                    type="radio"
                    :value="color"
                    v-model="selectedColor"
                    class="sr-only">
              </label>
            </div>
          </div>

          <div class="flex justify-end">
            <button type="submit" class="text-white bg-neutral-900 hover:bg-black dark:bg-black rounded px-6 py-2 dark:hover:bg-neutral-900 transition-colors duration-200">
              Conectar
            </button>
          </div>
        </form>
      </div>
    </div>
    <div class="h-screen flex bg-neutral-300 dark:bg-gray-900 bg-black" v-else>
      <div class="mx-auto flex justify flex-col bg-indigo-500 mt-10 mb-10 rounded-xl shadow-2xl w-full max-w-2xl overflow-auto dark:shadow-indigo-950 dark:border-0 dark:bg-gray-800 dark:text-white">
        <div class="bg-indigo-500 text-4xl pt-4 pb-4 pl-10 rounded-t-xl font-bold text-zinc-200 text-shadow-lg/60 dark:bg-gray-800 dark:text-shadow-lg/60 dark:text-shadow-indigo-500">Chatroom</div>
        <div class="overflow-auto flex flex-col bg-white dark:bg-gray-900 p-5 h-[60vh]" ref="containerRef"></div>
        <div class="flex-shrink-0 px-4 p-8 min-h-[100px] overflow-auto flex flex-col bg-indigo-500 dark:bg-gray-800">
          <form @submit.prevent="sendMessage" class="border rounded-lg p-2 border-stone-300 flex bg-white text-black">
            <input type="text"
                   autofocus
                   v-model="message"
                   placeholder="Escribe tu mensaje..."
                   class="border border-stone-400 outline-none focus:ring-0 px-4 py-2 w-full border-none flex-grow">
            <button type="submit" class="px-4 py-2 bg-black text-white rounded">Enviar</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>