import './style.css';
import './app.css';

import {Download} from '../wailsjs/go/main/Downloader';
import {EventsOn} from '../wailsjs/runtime/runtime'

let listElement = document.getElementById("m3u8-list");
listElement.focus();
// let outputFileElement = document.getElementById("outputFileName")
let logElement = document.getElementById("logs");

EventsOn("log-writer", (e) => {
    console.log("e", e)
    logElement.innerText += e
})
// Setup the greet function
window.download = function () {
    logElement.innerText = "";
    let list = listElement.value;
    // Check if the input is empty
    if (list === "") return;
    try {
        Download(list)
            .then((result) => {
                // Update result with data back from App.Greet()
                logElement.innerText += "downloaded to = " + result;
            })
            .catch((err) => {
                logElement.innerText += err;
            });
    } catch (err) {
        logElement.innerText += err;
    }
};
