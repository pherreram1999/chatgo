<script lang="ts" setup>

import {ref} from "vue";


enum Steps {
    Login,
    Chat
}

interface MessageResponse {
    usuario: string,
    cuerpo: string
}

const host = 'ws://localhost:8080/ws'
const username = ref<string>()
const message = ref<string>()


const step = ref<Steps>(Steps.Login)

let connection: WebSocket

const containerRef = ref()

const newNode = (tag: string) => document.createElement(tag)

const insertMessageNode = (msg: string) => {
    if(!containerRef.value) return
    const messageBody: MessageResponse = JSON.parse(msg)

    const divParent = newNode('div')
    divParent.className = 'px-4 py-2 border-b border-stone-300'
    const userLbl = newNode('span')
    userLbl.innerText = messageBody.usuario
    userLbl.className = 'block font-semibold'
    const bodyText = newNode('p')
    bodyText.innerText = messageBody.cuerpo
    divParent.appendChild(userLbl)
    divParent.appendChild(bodyText)

    containerRef.value.appendChild(divParent)

}

const handleRecieveMessage = (event: MessageEvent) => {
    insertMessageNode(event.data)
}

const sendMessage = () => {
    if(!connection) return
    connection.send(
        JSON.stringify({usuario: username.value, cuerpo: message.value})
    )
    message.value = ''
}

const conectar = () => {
    connection = new WebSocket(host)
    connection.onmessage = handleRecieveMessage
    step.value = Steps.Chat
}



</script>

<template>
    <div>
        <div v-if="step === Steps.Login" class="h-screen flex justify-center items-center">
            <div>
                <form @submit.prevent="conectar" class="p-4 rounded shadow-md ronded-lg border border-stone-200 w-[24rem]">
                    <h1 class="text-4xl font-semibold text-center mb-4">WebChat</h1>

                    <input type="text"
                           name="username"
                           id="username"
                           v-model="username"
                           autofocus
                           required
                           placeholder="Nombre de usuario"
                           class="bg-stone-100 w-full rounded px-4 py-2 border-stone-500">
                    <div class="flex justify-end mt-4">
                        <button type="submit" class="text-white bg-black rounded px-4 py-2">
                            Conectar
                        </button>
                    </div>
                </form>
            </div>
        </div>
        <div class="h-screen flex" v-else>
            <div class="w-[60rem] h-full mx-auto flex justify flex-col">
                <div class="flex-grow" ref="containerRef"></div>
                <div class="flex-shrink-0 px-4 mb-[4rem]">
                    <form @submit.prevent="sendMessage" class="border rounded-lg p-2 border-stone-300 flex">
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