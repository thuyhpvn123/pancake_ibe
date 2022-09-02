
let socket = new WebSocket("ws://127.0.0.1:3000/ws");
        console.log("Attempting Connection...");

        // socket.onopen = () => {
        //     console.log("Successfully Connected");
        //     socket.send("Hi From the Client!")
        // };
        
        socket.onclose = event => {
            console.log("Socket Closed Connection: ", event);
            socket.send("Client Closed!")
        };

        socket.onerror = error => {
            console.log("Socket Error: ", error);
        };
        const input = document.getElementById("input");
        const output = document.getElementById("output");
        // const socket = new WebSocket("ws://localhost:8000/echo");
        socket.onopen = function () {
        output.innerHTML += "Status: Connected\n";
            };
        // receive message from server and write on html output
        // socket.onmessage = function (e) {         
        // output.innerHTML += "Server: " + e.data + "\n";
        //     };
        //send message to server a
        function send() {
        socket.send(input.value);
        input.value = "";
            }
        socket.addEventListener('message', (event) => {
            writeToScreen('<span style = "color: blue;">Server: ' +
            event.data+'</span>'); 
            // socket.close();
            console.log('Message from server ', event.data);
        });
        function writeToScreen(message) {
         var pre = document.createElement("p"); 
         pre.style.wordWrap = "break-word"; 
         pre.innerHTML = message; 
         output.appendChild(pre);
      }
      const wrapper = document.querySelector(".wrapper"),
      qrInput = wrapper.querySelector(".form input"),
      generateBtn = wrapper.querySelector(".form button"),
      qrImg = wrapper.querySelector(".qr-code img");
      let preValue;
      
      generateBtn.addEventListener("click", () => {
          let qrValue = qrInput.value.trim();
          if(!qrValue || preValue === qrValue) return;
          preValue = qrValue;
          generateBtn.innerText = "Generating QR Code...";
          qrImg.src = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${qrValue}`;
          qrImg.addEventListener("load", () => {
              wrapper.classList.add("active");
              generateBtn.innerText = "Generate QR Code";
          });
      });
      
      qrInput.addEventListener("keyup", () => {
          if(!qrInput.value.trim()) {
              wrapper.classList.remove("active");
              preValue = "";
          }
      });