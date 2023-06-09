//import onMount from "svelte"
let time = ""
let para = document.getElementById("time")
const evtSrc = new EventSource("http://localhost:3500/event")
/*const messages = await fetch("http://localhost:3500/getMessages").then(
    (response) =>{
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`)
        }
        return response.text();
    }
).then((text) => para.textContent = text).catch((error) => para.textContent = `Could not fetch verse: ${error}`);

if(messages.status !== 200){
    console.log("Could not to the server")
}
*/
evtSrc.onmessage = function (event){
    time = event.data
    let para = document.getElementById("time")
    let node = document.createTextNode(time)
    para.appendChild(node)
    console.log(time)
}

evtSrc.onerror =  function (event){
    console.log(event)
    console.log("ошибка")
}

/*
async function getTime() {
    const res = await fetch("http://localhost:3500/time")
    if(res.status !== 200){
        console.log("Could not to the server")
    }
    console.log('test')
    return document.getElementById("time").innerText = "Time: " + time
}
getTime()

 */
/*
onMount( () => {
    const evtSrc = new EventSource("http://localhost:3500/event")
    evtSrc.onmessage = function (event){
        time = event.data
    }

    evtSrc.onerror =  function (event){
        console.log(event)
    }
})

 */
/*
const source = new EventSource("http://localhost:3000/sse")
source.onmessage = (event) => {
    console.log("OnMessage Called:")
    console.log(event)
    console.log(JSON.parse(event.data))
}*/