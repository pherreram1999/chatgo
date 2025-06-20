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

const host = 'ws://192.168.8.3:8081/ws'
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
    const uri = new URL(host)
    uri.searchParams.append('username',username.value)
    connection = new WebSocket(uri.toString())
    connection.onmessage = handleRecieveMessage
    step.value = Steps.Chat
}



</script>

<template>
    <div>
        <div v-if="step === Steps.Login" class="h-screen flex justify-center items-center bg-neutral-300 dark:bg-gray-900">
            <div>
                <form @submit.prevent="conectar" class="bg-white p-6 rounded shadow-2xl text-shadow-black text-shadow-lg/60 rounded-2xl border border-stone-200 w-[24rem] dark:shadow-indigo-950 dark:border-0 dark:bg-gray-800 dark:text-white dark:text-shadow-indigo-500 dark:text-shadow-lg/60">
                    <h1 class="text-4xl font-semibold text-center mb-10">WebChat Express</h1>

                    <input type="text"
                           name="username"
                           id="username"
                           v-model="username"
                           autofocus
                           required
                           placeholder="Nombre de usuario"
                           class="bg-stone-100 w-full rounded px-4 py-2 border-stone-500 text-black focus:outline-none focus:ring-2 focus:ring-black dark:focus:ring-indigo-500 focus:border-transparent placeholder-gray-400">
                    <div class="flex justify-end mt-4">
                        <button type="submit" class="text-white bg-neutral-900 hover:bg-black dark:bg-black rounded px-4 py-2 dark:hover:bg-neutral-900">
                            Conectar
                        </button>
                    </div>
                </form>
            </div>
        </div>
        <div class="h-screen flex bg-neutral-300 dark:bg-gray-900 bg-black" v-else>
            <div class="mx-auto flex justify flex-col bg-white mt-10 mb-10 rounded-xl shadow-2xl w-full max-w-2xl overflow-auto bg-white dark:shadow-indigo-950 dark:border-0 dark:bg-gray-800 dark:text-white">
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