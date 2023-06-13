const messageInput = document.getElementById('message-input');
const sendButton = document.getElementById('send-button');
const messageList = document.getElementById('messages');

function addMessage(message){
    if (message) {
        const messageElement = document.createElement('li');
        // messageElement.textContent = message;
        // messageElement.setAttribute("class","message-item");
        messageElement.innerHTML=message;
        messageList.appendChild(messageElement);
        messageInput.value = '';
    }
}


sendButton.addEventListener('click', () => {
    const message = messageInput.value.trim();
    addMessage(message);
    $.ajax({
        method:"POST",
        url:"/chat",
        contentType:"application/json; charset=utf-8",
        data:JSON.stringify({
            req:message   
       }),
        success:function(resp){
            console.log(resp);
            addMessage(resp.msg);
        },
        error:function(){
            addMessage("出错了，要不等会试试");
        }
    })
    
});