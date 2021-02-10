// 用户信息map，键为用户唯一标识，值为用户信息
// todo: ids并发访问处理
var ids = new Map();
var myName = "";
var myId = "";
// 点击准备/取消准备按钮后，按钮上显示的文字
var readyChangeStatus = {'ready': '取消准备', 'cancel': '准备'};
var occupations = ["scientist", "engineer", "doctor", "driver"];
// 发送消息文本框中的内容
var sendMessage = function (ws) {
    let msg = $("#msg").val();
    if (msg.length > 0) {
        let messageInfo = { id: myId, msg: msg, username: myName };
        ws.send(JSON.stringify(messageInfo));
        $("#your-message").append($("<li>").text(msg).attr("class", "msg"));
        // 最多显示3条最新消息
        if ($("#your-message").children().length > 3) {
            $("#your-message").children().eq(0).remove();
        }
        $("#msg").val('');
    } else {
        alert("消息不能为空！");
    }
};
// 刷新玩家信息
// 参数：{"name1": "uuid1", "name2": "uuid2", ...}
var updateUsers = function(users) {
    // console.log(users);
    // 将参数的键/值对调，因为uuid是唯一的，但用户名可能与已有id冲突
    let newIds = {};
    $.each(users, function(i, e) {
        newIds[e] = i;
    });
    let userNum = Object.keys(newIds).length;
    $("#player-num").text(userNum);
    // 删除不在线的用户
    $(".content").children().each(function (i, e) {
        if (!newIds.hasOwnProperty($(e).attr('id'))) {
            console.log("remove", $(e).attr('id'));
            ids.delete($(e).attr('id'));
            $(e).remove();
        }
    });
    // 如果是新增用户，添加box
    for (let id in newIds) {
        if (id === myId) {
            continue;
        }
        // 当id不在当前用户id map中时，进行插入
        if (!ids.has(id)) {
            ids.set(id, true);
            let userBoxDiv = $('<div>').attr('class', 'box').attr('id', id);
            // 添加用户名
            let userNameEle = $('<h2>').text(newIds[id]);
            userBoxDiv.append(userNameEle);
            // 添加道具栏
            let toolBoxDiv = $('<div>').attr('class', 'toolbox');
            let toolBoxHead = $('<h2>').text('道具栏');
            let toolsDiv = $('<div>');
            toolBoxDiv.append(toolBoxHead);
            toolBoxDiv.append(toolsDiv);
            userBoxDiv.append(toolBoxDiv);
            $(".content").append(userBoxDiv);
        }
    }
};
// 添加其他用户的消息
var appendMessage = function(req) {
    if (req["id"] === myId) {
        return;
    }
    // 不是自己发送的消息，将其添加到相应用户的框中
    let eleId = "#" + req["id"];
    // console.log("message from", eleId);
    $('<li>').text(req["msg"]).attr("class", "msg").appendTo($(eleId));
    // 最多显示3条最新消息，消息前存在用户信息和道具栏，所以是5不是3
    if ($(eleId).children().length > 5) {
        $(eleId).children().eq(2).remove();
    }
};
// 定时向服务器发送POST请求进行身份信息校验
var checkUser = function(callback) {
    // callback函数只执行一次，之后如果身份校验成功什么都不做，校验失败直接返回
    let userInfo = { username: myName, id: myId };
    let reqData = JSON.stringify(userInfo);
    $.post('/checkuser', reqData, function (data, status) {
        if (status == "success") {
            // console.log("checkuser:", data);
            if (data == "ok") {
                callback();
            } else {
                $(window).attr("location", "/");
            }
        }
    });
    // 定时校验id和name，解决多开窗口后，原窗口“假在线”问题
    window.setInterval(function() {
        $.post('/checkuser', reqData, function (data, status) {
            if (status == "success") {
                // console.log("checkuser:", data);
                if (data != "ok") {
                    $(window).attr("location", "/");
                }
            }
        });
    }, 1000);
};
// 玩家准备/取消准备
var userChangeReady = function(action) {
    $.get("/gamectrl/" + action + "/" + myName, function(data, status) {
        if (status == "success") {
            if (data == "ok") {
                $("#btn-ready").text(readyChangeStatus[action]);
            } else {
                alert("游戏正在进行中，请稍等...");
            }
        }
    });
};
// 游戏开始、结束
var changeGameStatus = function(running) {
    console.log(running);
    if (running == "yes") {
        $("#game").show(1000);
    } else {
        $("#game").hide(1000);
    }
};
// 更新游戏信息
var updateGame = function(req) {
    console.log(req);
    // 更新回合数
    $("#gameround").text(" Round " + req["round"]);
    // 更新地图显示
    let map = req["map"];
    for (let i = 0; i < 7; i++) {
        for (let j = 0; j < 7; j++) {
            // 添加病毒
            if (map[i][j]["virus"] != 0) {
                let virusImg = $('<img>').attr('src', '/static/img/virus.png').attr('class', 'level'+map[i][j]["virus"]);
                $("#" + i + "-" + j).children().eq(0).append(virusImg);
            }
            // 添加研究所
            if (map[i][j]["institute"] != 0) {
                let instituteImg = $('<img>').attr('src', '/static/img/cap1.png');
                $("#" + i + "-" + j).children().eq(1).append(instituteImg);
            }
            // 添加玩家位置
            $("#" + i + "-" + j).children().eq(2).text(map[i][j]["player"]);
        }
    }
    // 更新玩家职业
    let players = req["players"];
    $.each(players, function(i, e) {
        let name = $("#" + e["posx"] + "-" + e["posy"]).children().eq(2).text();
        $("#" + e["posx"] + "-" + e["posy"]).children().eq(2).text(name + "(" + occupations[e["occupation"]] + ")");
    });
}
// WebSocket处理函数
var wsProcess = function () {
    $("#username").text(myName);
    let ws = new WebSocket('ws://' + window.location.host + '/ws');
    // 建立连接后发送空消息带id，建立用户和连接的关系，断开连接时清除登录状态
    ws.onopen = function (e) {
        // console.log("websocket opened!");
        let messageInfo = { id: myId, msg: "", username: myName };
        ws.send(JSON.stringify(messageInfo));
    }
    // 服务器WebSocket消息处理
    ws.onmessage = function (e) {
        // console.log(e.data);
        // todo: e.data类型判断
        let reqData = JSON.parse(e.data);
        if (reqData["users"] !== undefined) { // 收到服务器定时发送的用户信息等
            updateUsers(reqData["users"]);
        } else if (reqData["id"] !== undefined) { // 收到消息
            appendMessage(reqData);
        } else if (reqData["running"] !== undefined) { // 开始、结束游戏
            changeGameStatus(reqData["running"]);
        } else if (reqData["map"] !== undefined) { // 处理地图信息
            updateGame(reqData);
        } else {
            // 处理其他类型消息
        }
    };
    // 每次发消息时在服务器校验id
    $("#msg").on("keydown", function (event) {
        // 回车键发送消息
        var keyCode = event.keyCode || event.which;
        if (keyCode == 13) {
            // event.preventDefault();
            sendMessage(ws);
        }
    });
    $("#btn-send").on("click", function () {
        sendMessage(ws);
    });
    $("#btn-ready").on("click", function() {
        if ($(this).text() == "准备") {
            userChangeReady('ready');
        } else {
            userChangeReady('cancel');
        }
    });
};
$(function () {
    // name和id只在初始化时获取一次，如果同一主机又登录一个相同用户，校验会失败，防止多开窗口
    myName = getCookie("username");
    myId = getCookie("id");
    // console.log("My Id:", myId);
    // console.log("My Name:", myName);
    // 用户鉴权，如果鉴权成功，处理服务器WebSocket数据
    if (myName != "" && myId != "") {
        checkUser(wsProcess);
    } else {
        $(window).attr("location", "/");
    }
});