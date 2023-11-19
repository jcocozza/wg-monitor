function toggleClass(oldClass ,newClass, elementId) {
    var element = document.getElementById(elementId);
    if (element.classList.contains(oldClass)) {
        element.classList.remove(oldClass);
        element.classList.add(newClass); // Replace 'new-class' with the class you want to toggle to
    } else {
        element.classList.remove(newClass); // Replace 'new-class' with the class you want to toggle to
        element.classList.add(oldClass);
    }
}

function updateInterface(data){
    console.log(data)
    data.forEach(peer => {
        console.log(peer.publicKey)
        var peerId = peer.publicKey;
        var publickey = document.getElementById(peerId + "-publickey");
        var endpoint = document.getElementById(peerId + "-endpoint");
        var allowedips = document.getElementById(peerId + "-allowedips");
        var latesthandshake = document.getElementById(peerId + "-latesthandshake");
        var sent = document.getElementById(peerId + "-sent");
        var received = document.getElementById(peerId + "-received");

        var loadingPeer = document.getElementById(peerId + "-index")
        var peerStatus = document.getElementById(peerId + "-status")
        if (peer.info.status) {
            loadingPeer.setAttribute("aria-busy", "false");
            peerStatus.innerHTML = '<span class="statusDot online"></span>';
    
        } else {
            loadingPeer.setAttribute("aria-busy", "false");
            peerStatus.innerHTML = '<span class="statusDot offline"></span>';
        }

        publickey.innerHTML = peer.publicKey;     
        endpoint.innerHTML = peer.info.endPoint; 
        allowedips.innerHTML = peer.allowedIPs; 
        latesthandshake.innerHTML = peer.info.latestHandshake; 
        sent.textContent = peer.info.transfer.Sent; 
        received.textContent = peer.info.transfer.Received; 

    })
}
function fetchUpdateInterface(confName) {
    fetch("/api/update/configurations/" + confName)
        .then(response => response.json())
        .then(data => updateInterface(data))
        .catch(error => console.error("Error:", error));
}

function closeElement(formId) {
    document.getElementById(formId).style.display = 'none';
}


function fetchUpdateAllConfigurations(){
    fetch("/api/update/configurations/all")
    .then(response => response.json())
    .then(data => UpdateAllConfigurations(data))
    .catch(error => console.error("Error:", error))
}
setInterval(fetchUpdateAllConfigurations, 5000)

function UpdateAllConfigurations(data){
    console.log(data)

    for (let confName in data) {
        let cStatus = data[confName];
        confStatus = document.getElementById(confName+"-status")

        if (cStatus) {
            //loadingPeer.setAttribute("aria-busy", "false");
            confStatus.innerHTML = '<span class="statusDot online"></span>';
    
        } else {
            //loadingPeer.setAttribute("aria-busy", "false");
            confStatus.innerHTML = '<span class="statusDot offline"></span>';
        }
    }
}
