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
    data.forEach(peer => {
        var peerId = peer.PublicKey;
        var publickey = document.getElementById(peerId + "-publickey");
        var endpoint = document.getElementById(peerId + "-endpoint");
        var allowedips = document.getElementById(peerId + "-allowedips");
        var latesthandshake = document.getElementById(peerId + "-latesthandshake");
        var sent = document.getElementById(peerId + "-sent");
        var received = document.getElementById(peerId + "-received");

        var loadingPeer = document.getElementById(peerId + "-index")
        var peerStatus = document.getElementById(peerId + "-status")
        if (peer.Info.Status) {
            loadingPeer.setAttribute("aria-busy", "false");
            peerStatus.innerHTML = '<span class="statusDot online"></span>';
    
        } else {
            loadingPeer.setAttribute("aria-busy", "false");
            peerStatus.innerHTML = '<span class="statusDot offline"></span>';
        }

        publickey.innerHTML = peer.PublicKey;     
        endpoint.innerHTML = peer.Info.EndPoint; 
        allowedips.innerHTML = peer.AllowedIPs; 
        latesthandshake.innerHTML = peer.Info.LatestHandshake; 
        sent.innerHTML = peer.Info.Transfer.Sent; 
        received.innerHTML = peer.Info.Transfer.Received; 

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
