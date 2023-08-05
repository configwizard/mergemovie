import './style.css';
import './app.css';

import {DiscoverM3u8MasterLinks, RetrieveVariants, Download, DirectDownload, OpenInDefaultBrowser} from '../wailsjs/go/main/Downloader';
import {EventsOn} from '../wailsjs/runtime/runtime'

let listElement = document.getElementById("m3u8-list");
listElement.focus();
// let outputFileElement = document.getElementById("outputFileName")
const logElement = document.getElementById('logs');
logElement.scrollTop = logElement.scrollHeight;
window.onload = function() {
    logElement.innerText = 'ready...';
};

EventsOn("log-writer", (e) => {
    console.log("e", e)
    logElement.innerText += e
})
async function isM3U8(url) {
    //CORS issue here
    // const response = await fetch(url, { method: 'HEAD' });
    // const contentType = response.headers.get('Content-Type');
    // return contentType === 'application/vnd.apple.mpegurl' || contentType === 'application/x-mpegURL';
    return url.endsWith('.m3u8');
}
function splitURL(fullUrl) {
    const urlObject = new URL(fullUrl);
    const pathParts = urlObject.pathname.split('/');
    const variantPath = pathParts.pop(); // Get the last part of the path
    urlObject.pathname = pathParts.join('/'); // Join the remaining parts back together
    const masterUrl = urlObject.toString(); // Get the master URL without the variant path
    return { masterUrl, variantPath };
}

// Setup the greet function
window.download = async function () {
    logElement.innerText = "";
    let lnk = listElement.value;
    // Check if the input is empty
    if (lnk === "") return;
    try {
        if (await isM3U8(lnk)) {
            try {
                DirectDownload(lnk)
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

            return
        }
        logElement.innerText += "detecting m3u8 links";
        DiscoverM3u8MasterLinks(lnk)
            .then((result) => {
                const notificationsDiv = document.getElementById('notifications');
                notificationsDiv.innerText = '';
                if (result == null || result.length == 0) {
                    logElement.innerText += "no m3u8 videos found"
                    return
                }
                result.forEach((url) => {
                    const notificationDiv = document.createElement('div');
                    notificationDiv.className = 'notification';

                    const textNode = document.createTextNode(`Found ${url} - do you want to search this for versions?`);
                    notificationDiv.appendChild(textNode);

                    const button = document.createElement('button');
                    button.innerText = 'Select';
                    button.onclick = function() {
                        // Handle the selection here, you may call a function that handles the selected URL
                        handleSelectedURL(url, notificationDiv);
                    };

                    notificationDiv.appendChild(button);

                    notificationsDiv.appendChild(notificationDiv);
                });
                // Update result with data back from App.Greet()
                logElement.innerText += "discovered " + result;
            })
            .catch((err) => {
                logElement.innerText += err;
            });
    } catch (err) {
        logElement.innerText += err;
    }
};

function handleSelectedURL(url, parentNotificationDiv) {
    RetrieveVariants(url).then((variants) => {
        displayVariants(parentNotificationDiv, url, variants);
    })
}
function displayVariants(parentNotificationDiv, masterUrl, variants) {
    const variantsDiv = document.createElement('div');
    variantsDiv.className = 'variants-container';
    variantsDiv.className = 'variants';
    parentNotificationDiv.innerHTML = '';
    variants.forEach((variant) => {
        logElement.innerText += "found variant " + variant;
        const quality = variant.match(/\d+/)[0]; // Extract quality from the name
        const variantDiv = document.createElement('div');
        variantDiv.className = 'variant';

        const textNode = document.createTextNode(`Choose quality ${quality}p?`);
        variantDiv.appendChild(textNode);

        const button = document.createElement('button');
        button.innerText = 'Select';
        button.onclick = function() {
            // Handle the selection here
            handleSelectedVariant(masterUrl, variant);
        };

        variantDiv.appendChild(button);

        variantsDiv.appendChild(variantDiv);
    });

    parentNotificationDiv.appendChild(variantsDiv);
}

function handleSelectedVariant(masterUrl, variant) {
    // Handle the selected variant here
    console.log(`Selected Variant: ${masterUrl} - ${variant}`);
    Download(masterUrl, variant).then(result => {
        logElement.innerText += "downloaded " + result;

        // Create success notification
        const successDiv = document.createElement('div');
        successDiv.className = 'success-notification';
        successDiv.innerText = 'Download successful!';

        // Append success notification to a specific area (e.g., notifications area)
        const notificationsDiv = document.getElementById('notifications');
        notificationsDiv.appendChild(successDiv);
    }).catch((err) => {
        logElement.innerText += err;
    });
}


window.openInDefaultBrowser = async(txt) => {
    try {
        OpenInDefaultBrowser(txt)
            .then((result) => {
            })
            .catch((err) => {
                logElement.innerText += err;
            });
    } catch (err) {
        logElement.innerText += err;
    }
}