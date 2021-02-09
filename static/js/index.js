var sendMessage = function(ws, myId, myUsername) {
    let msg = $("#msg").val();
    if (msg.length > 0) {
        let messageInfo = { id: myId, msg: msg, username: myUsername };
        ws.send(JSON.stringify(messageInfo));
        $("#your-message").append($("<li>").text(msg).attr("class", "msg"));
        if ($("#your-message").children().length > 3) {
            $("#your-message").children().eq(0).remove();
        }
        $("#msg").val('');
    } else {
        alert("消息不能为空！");
    }
}
var wsProcess = function (myId, myUsername) {
    $("#username").text(myUsername);
    let ws = new WebSocket('ws://' + window.location.host + '/ws');
    let ids = new Map();
    // 建立连接后发送空消息带id，建立用户和连接的关系，断开连接时清除登录状态
    ws.onopen = function (e) {
        console.log("websocket opened!");
        let messageInfo = { id: myId, msg: "", username: myUsername };
        ws.send(JSON.stringify(messageInfo));
    }
    // 服务器WebSocket消息处理
    ws.onmessage = function (e) {
        // console.log(e.data);
        // todo: e.data类型判断
        let message = JSON.parse(e.data);
        if (message["users"] !== undefined) { // 收到服务器定时发送的用户状态信息等
            let users = message["users"];
            // console.log(users);
            let userNum = Object.keys(users).length;
            $("#player-num").text(userNum);
            // 删除不在线的用户
            $(".content").children().each(function (i, e) {
                // console.log($(e).attr('id'));
                if (!users.hasOwnProperty($(e).attr('id'))) {
                    console.log("remove", $(e).attr('id'));
                    ids.delete($(e).attr('id'));
                    $(e).remove();
                }
            });
            // $(".content").text("");
            // message: {users: [id1: name1, id2: name2...]}
            // 如果是新增用户，添加box
            for (let id in users) {
                if (id === myId) {
                    continue;
                }
                if (!ids.has(id)) {
                    ids.set(id, true);
                    let userBoxDiv = $('<div>').attr('class', 'box').attr('id', id);
                    let userNameEle = $('<h2>').text(users[id]);
                    userBoxDiv.append(userNameEle);
                    let toolBoxDiv = $('<div>').attr('class', 'toolbox');
                    let toolBoxHead = $('<h2>').text('道具栏');
                    let toolsDiv = $('<div>');
                    toolBoxDiv.append(toolBoxHead);
                    toolBoxDiv.append(toolsDiv);
                    userBoxDiv.append(toolBoxDiv);
                    $(".content").append(userBoxDiv);
                }
            }
        } else { // 收到消息
            if (message["id"] !== myId) {
                let eleId = "#" + message["id"];
                console.log("message from", eleId);
                $('<li>').text(message["msg"]).attr("class", "msg").appendTo($(eleId));
                if ($(eleId).children().length > 5) {
                    $(eleId).children().eq(2).remove();
                }
            }
        }
    };
    // 每次发消息时在服务器校验id
    $("#msg").on("keydown", function (event) {
        var keyCode = event.keyCode || event.which;
        console.log(typeof keyCode);
        if (keyCode == 13) {
            // event.preventDefault();
            sendMessage(ws, myId, myUsername);
        }
    });
    $("#btn-send").on("click", function () {
        sendMessage(ws, myId, myUsername);
    });
};
$(function () {
    let username = getCookie("username");
    let id = getCookie("token");
    console.log("id:", id);
    console.log("username:", username);
    let userInfo = { username: username, id: id };
    if (username != "" && id != "") {
        // 用户鉴权
        $.post('/checkuser', JSON.stringify(userInfo), function (data, status) {
            if (status == "success") {
                console.log("checkuser:", data);
                if (data == "ok") {
                    wsProcess(id, username);
                } else {
                    $(window).attr("location", "/");
                }
            }
        });
        // todo: 定时检查cookie，问题：多开窗口检查

    } else {
        $(window).attr("location", "/");
    }
});