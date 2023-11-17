function updateInterfaces(data) { 
    data.forEach(iface => {
        var interfaceId = "interface-" + iface.Name;
        iface["Peers"].forEach(peer => {
            var peerId = "peer-" + peer.PublicKey;
            var publickey = document.getElementById(peerId + "-publickey");
            var endpoint = document.getElementById(peerId + "-endpoint");
            var allowedips = document.getElementById(peerId + "-allowedips");
            var latesthandshake = document.getElementById(peerId + "-latesthandshake");
            var transfer = document.getElementById(peerId + "-transfer");

            publickey.innerHTML = peer.PublicKey;     
            endpoint.innerHTML = peer.EndPoint; 
            allowedips.innerHTML = peer.AllowedIPs; 
            latesthandshake.innerHTML = peer.LatestHandshake; 
            transfer.innerHTML = peer.Transfer; 
        })
    });
}

function getInterfaces() {
    fetch('/api/getInterfaces')
        .then(response => response.json())
        .then(data => updateInterfaces(data))
        .catch(error => console.error("Error:", error));
}

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


function newPopup(url) {
    popupWindow = window.open(
		url,'popUpWindow','height=300,width=400,left=10,top=10,resizable=yes,scrollbars=yes,toolbar=yes,menubar=no,location=no,directories=no,status=yes')
}

