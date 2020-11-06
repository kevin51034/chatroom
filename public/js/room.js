const chatForm = document.getElementById('chat-form');
const chatMessages = document.querySelector('.chat-messages');
const roomName = document.getElementById('room-name');
const userList = document.getElementById('users');

// Get username and room from URL
const { username, room } = Qs.parse(location.search, {
  ignoreQueryPrefix: true
});

console.log(username)
console.log(room)

const socket = io();


/*
// Join chatroom
socket.emit('joinRoom', { username, room });

// Get room and users
socket.on('roomUsers', ({ room, users }) => {
  outputRoomName(room);
  outputUsers(users);
});


// Message from server
socket.on('message', message => {
  //console.log(message);
  outputMessage(message);

  // Scroll down
  chatMessages.scrollTop = chatMessages.scrollHeight;
});
*/

// *****
if (window["WebSocket"]) {
  conn = new WebSocket("ws://" + document.location.host + "/ws" + `?username=${username}&room=${room}`);
  console.log("build connection")
  outputRoomName(room);
  outputUsers(users);
  socket.emit('joinRoom', { username, room });

  conn.onclose = function (evt) {
            /*var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);*/
            console.log("Connection close")
  };
  conn.onmessage = function (evt) {
      var messages = evt.data.split('\n');
      /*
      for (var i = 0; i < messages.length; i++) {
          var item = document.createElement("div");
          item.innerText = messages[i];
          appendLog(item);
      }*/
      outputMessage(messages)
      console.log(messages)
  };
} else {
  /*var item = document.createElement("div");
  item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
  appendLog(item);*/
  console.log("Your browser does not support WebSockets.")
}
// *****

// Message submit
chatForm.addEventListener('submit', e => {
  e.preventDefault();
  // Get message text
  let msg = e.target.elements.msg.value;
  console.log(msg)
  conn.send(msg);

  msg = msg.trim();
  
  if (!msg){
    return false;
  }
  //outputMessage(msg)
  /*
  // Emit message to server
  socket.emit('chatMessage', msg);
  */

  // Clear input
  e.target.elements.msg.value = '';
  e.target.elements.msg.focus();
});

// Output message to DOM
function outputMessage(message) {
  const div = document.createElement('div');
  div.classList.add('message');
  const p = document.createElement('p');
  p.classList.add('meta');

  //p.innerText = message.username;
  p.innerText = username;

  p.innerHTML += `<span>${message.time}</span>`;
  div.appendChild(p);
  const para = document.createElement('p');
  para.classList.add('text');
  //para.innerText = message.text;
  para.innerText = message;
  div.appendChild(para);
  document.querySelector('.chat-messages').appendChild(div);
}

// Add room name to DOM
function outputRoomName(room) {
  roomName.innerText = room;
}

// Add users to DOM
function outputUsers(users) {
  userList.innerHTML = '';
  /*users.forEach(user=>{
    const li = document.createElement('li');
    li.innerText = user.username;
    userList.appendChild(li);
  });*/
  const li = document.createElement('li');
  li.innerText = username;
  userList.appendChild(li);
}
