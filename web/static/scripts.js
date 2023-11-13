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


function updatePeerMetaStatus(data){
    data["Peers"].forEach(peer => {
        peerMetaStatusId = peer.PublicKey + "-MetaStatus"
        var peerMetaStatus = document.getElementById(peerMetaStatusId);

        var peerId = peer.PublicKey;
        var publickey = document.getElementById(peerId + "-publickey");
        var endpoint = document.getElementById(peerId + "-endpoint");
        var allowedips = document.getElementById(peerId + "-allowedips");
        var latesthandshake = document.getElementById(peerId + "-latesthandshake");
        var sent = document.getElementById(peerId + "-sent");
        var recieved = document.getElementById(peerId + "-recieved");

        var loadingPeer = document.getElementById(peerId + "-index")

        if (peer.MetaStatus) {
            loadingPeer.setAttribute("aria-busy", "false");
            peerMetaStatus.innerHTML = '<span class="statusDot online"></span>';
    
        } else {
            loadingPeer.setAttribute("aria-busy", "false");
            peerMetaStatus.innerHTML = '<span class="statusDot offline"></span>';
        }

        publickey.innerHTML = peer.PublicKey;     
        endpoint.innerHTML = peer.EndPoint; 
        allowedips.innerHTML = peer.AllowedIPs; 
        latesthandshake.innerHTML = peer.LatestHandshake; 
        sent.innerHTML = peer.Sent; 
        recieved.innerHTML = peer.Recieved; 

    })
}
function getPeerMetaStatus(interfaceName) {
    fetch("/api/configurations/" + interfaceName)
        .then(response => response.json())
        .then(data => updatePeerMetaStatus(data))
        .catch(error => console.error("Error:", error));
}

