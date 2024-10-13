// Create a new room flow
document.getElementById('create-btn').addEventListener('click', showCreateRoomScreen);
document.getElementById('create-room-btn').addEventListener('click', createRoom);

// Join an existing room flow
document.getElementById('join-btn').addEventListener('click', showJoinRoomScreen);
document.getElementById('join-room-btn').addEventListener('click', joinRoom);

document.getElementById('mute-btn').addEventListener('click', toggleMute);
document.getElementById('video-btn').addEventListener('click', toggleVideo);

let localStream;
let peerConnection;
let isMuted = false;
let isVideoStopped = false;

let ws;

async function showCreateRoomScreen() {
  console.log('Creating room');
  document.getElementById('welcome-screen').style.display = 'none';
  document.getElementById('enter-room').style.display = 'block';
  document.getElementById('create-screen').style.display = 'block';
}

async function createRoom() {
  console.log('Creating new conference room');
  
  // Create a new room
  const resp = await fetch('/create')
  const { room_id } = await resp.json();
  console.log('Room created:', room_id);
  
  await sendJoinRoom(room_id);
}

async function showJoinRoomScreen() {
  console.log('Joining room');
  document.getElementById('welcome-screen').style.display = 'none';
  document.getElementById('enter-room').style.display = 'block';
  document.getElementById('join-screen').style.display = 'block';

}

async function joinRoom() {
  console.log('Joining room');
  const roomID = document.getElementById('room-id').value;
  if (!roomID) {
    alert('Please enter the room ID');
    return;
  }
  await sendJoinRoom(roomID);
}

async function sendJoinRoom(roomID) {
  const name = document.getElementById('name').value;
  if (!name) {
      alert('Please enter your name');
      return;
  }
  console.log('Joining room:', roomID);
  await activateStream();

  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';

  ws = new WebSocket(`${protocol}://${window.location.host}/join?roomID=${roomID}`);

  ws.onopen = () => {
    console.log('Connected to the signaling server');
    ws.send(JSON.stringify({ type: 'join' }));

    // Show the participant view
    document.getElementById('enter-room').style.display = 'none';
    document.getElementById('participant-view').style.display = 'block';

  }

  ws.onmessage = async (message) => {
    const data = JSON.parse(message.data);
    switch (data.type) {
      case 'join':
        // call user
        callUser();
        break;
      case 'offer':
        // handle offer
        handleOffer(data.offer);
        break;
      case 'answer':
        // handle answer
        console.log("Receiving Answer")
        peerConnection.setRemoteDescription(
          new RTCSessionDescription(data.answer)
        );
        break;
      case 'iceCandidate':
        // handle ICE
        console.log("Receiving ICE Candidate")
        try {
          peerConnection.addIceCandidate(data.iceCandidate);
        } catch (e) {
          console.error('Error adding received ice candidate', e);
        }
        break;
      default:
        console.log(`Unknown message type: ${JSON.stringify(data)}`); break;
    }
  }

}

async function activateStream() {
  localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
  console.log(`Activated Stream: ${JSON.stringify(localStream)}`);

  const localVideo = document.createElement('video');
  localVideo.srcObject = localStream;
  localVideo.autoplay = true;
  localVideo.muted = true;
  document.getElementById('videos').appendChild(localVideo);
}

async function handleOffer(offer) {
    console.log("Received Offer, Creating Answer");
    peerConnection = createPeer();
    
    // Set the recieved offer as the remote description
    peerConnection.setRemoteDescription(
        new RTCSessionDescription(offer)
    );

    // Add data to the local description
    console.log(`Adding Answer tracks: ${JSON.stringify(localStream)}`);
    localStream.getTracks().forEach(track => peerConnection.addTrack(track, localStream));

    const answer = await peerConnection.createAnswer();
    console.log(`Created Answer: ${JSON.stringify(answer)}`);
    await peerConnection.setLocalDescription(answer);
    console.log(`Set Local Description: ${JSON.stringify(peerConnection.localDescription)}`);

    ws.send(
      JSON.stringify({ type: "answer", answer: peerConnection.localDescription })
    );
};


async function callUser() {
  console.log("Calling Other User");
  peerConnection = createPeer();
  
  localStream.getTracks().forEach(track => peerConnection.addTrack(track, localStream));
}

function createPeer() {
  console.log("Creating Peer Connection");
  const peer = new RTCPeerConnection({
      iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
  });

  peer.onnegotiationneeded = handleNegotiationNeeded;
  peer.onicecandidate = handleIceCandidateEvent;
  peer.ontrack = handleTrackEvent;

  return peer;
}

async function handleNegotiationNeeded() {
  console.log("Creating Offer");

  try {
      const myOffer = await peerConnection.createOffer();
      await peerConnection.setLocalDescription(myOffer);

      await ws.send(
        JSON.stringify({ type: 'offer', offer: peerConnection.localDescription })
      );
  } catch (err) {
    console.error('Error creating offer', err);
  }
}     
      
async function handleIceCandidateEvent(e) {
    console.log("Found Ice Candidate");
    if (e.candidate) {
        console.log(e.candidate);
        ws.send(
          JSON.stringify({ type: "iceCandidate",  iceCandidate: e.candidate })
        );
    }
}


function handleTrackEvent(e) {
    console.log("Received Tracks");
    console.log(e.streams)
    //partnerVideo.srcObject = e.streams[0];
    addRemoteStream(e.streams[0]);
}

function toggleMute() {
    localStream.getAudioTracks().forEach(track => track.enabled = !track.enabled);
    isMuted = !isMuted;
    document.getElementById('mute-btn').textContent = isMuted ? 'Unmute' : 'Mute';
}

function toggleVideo() {
    localStream.getVideoTracks().forEach(track => track.enabled = !track.enabled);
    isVideoStopped = !isVideoStopped;
    document.getElementById('video-btn').textContent = isVideoStopped ? 'Start Video' : 'Stop Video';
}

function addRemoteStream(stream) {
    console.log("Adding Remote Stream");
    const remoteVideo = document.createElement('video');
    remoteVideo.srcObject = stream;
    remoteVideo.autoplay = true;
    document.getElementById('videos').appendChild(remoteVideo);
}


