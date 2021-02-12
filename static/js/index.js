// 用户信息map，键为用户唯一标识，值为用户信息
// todo: ids并发访问处理
var ids = new Map();
var myName = "";
var myId = "";
var myOccupation = "";
var players = [];
var mapInfo = {};
var toDrive = "";
var researchProgress = 0;
// 点击准备/取消准备按钮后，按钮上显示的文字
var readyChangeStatus = { 'ready': '取消准备', 'cancel': '准备' };
var occupations = ["scientist", "engineer", "doctor", "driver"];
var occupationsDetail = [
    "可以在有研究所的格子内进行研究，每次消耗2行动力，进度+5%",
    "可以在没有研究所的格子内建造研究所，每次消耗4行动力",
    "消灭病毒消耗行动力为1，其他玩家为1+病毒等级",
    "移动1格消耗0.5行动力，可以搭载1名在同一格的其他玩家到达另一位置，被搭载的玩家不消耗行动力"];
var gameStatus = ["游戏未开始", "游戏进行中"];
var cmds = {
    ADD_VIRUS: 0, // 添加病毒
    SUB_VIRUS: 1,           // 清除病毒
    MOV_PLAYER: 2,          // 移动玩家
    SUB_PLAYER: 3,          // 删除玩家，如断开连接时执行
    ADD_INSTITUE: 4,        // 建造研究所，工程师可以执行
    DO_RESEARCH: 5         // 进行研究，科学家可以进行
}
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
// users是一个map：{"name1": "uuid1", "name2": "uuid2", ...}
var updateUsers = function (req) {
    let users = req["users"];
    let status = req["gamestatus"];
    if ($("#game-status").text() != gameStatus[status]) {
        $("#game-status").text(gameStatus[status]);
    }
    if (req["players"] != null) {
        $("#players").text(req["players"].join(','));
    }
    // 将参数的键/值对调，因为uuid是唯一的，但用户名可能与已有id冲突
    let newIds = {};
    $.each(users, function (i, e) {
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
var appendMessage = function (req) {
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
var checkUser = function (callback) {
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
    window.setInterval(function () {
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
var userChangeReady = function (action) {
    $.get("/gamectrl/" + action + "/" + myName, function (data, status) {
        if (status == "success") {
            console.log("userChangeReady: ", data);
            if (data == "ok") {
                $("#btn-ready").text(readyChangeStatus[action]);
            } else {
                alert("游戏正在进行中，请稍等...");
            }
        }
    });
};
// 玩家回合结束
var userFinish = function () {
    $("#btn-finish").hide();
    $.get("/gamectrl/finish/" + myName, function (data, status) {
        if (status == "success") {
            console.log("userFinish: ", data);
            if (data == "ok") {
                // 如果在这里隐藏按钮，由于是异步的，所以会导致开始了新回合时按钮还是没有恢复显示，实际是新回合显示了但又隐藏了
            } else if (data == "warn") {
                alert("请等待其他玩家结束回合...");
            } else {
                $("#btn-finish").show();
                alert("您不在游戏中...");
            }
        }
    });
};
// 游戏开始、结束
var changeGameStatus = function (req) {
    if (req["running"] == "yes") {
        $("#game").show(1000);
    } else {
        $("#game").hide(1000);
        $("#btn-ready").text("准备");
        if (req["message"] == "win") {
            alert("游戏胜利!");
        } else {
            alert("游戏失败!");
        }
    }
};
// 点击格子后，检查路径是否可行，向服务器请求移动
var moveUser = function (square) {
    // console.log($(square).attr('id'));
    let player = getPlayerObj();
    console.log("[moveUser]player: ", player);
    if (player == null) {
        return;
    }
    let pos = $(square).attr('id').split('-');
    let tarx = pos[0];
    let tary = pos[1];
    let energy = Math.abs(player["posx"] - tarx) + Math.abs(player["posy"] - tary);
    if (occupations[player["occupation"]] == 'driver') {
        energy /= 2.0
    }
    if (player["energy"] - energy < 0.0) {
        console.log("energy is not enough");
        return;
    }
    // todo: 司机带人
    let movInfo = { username: myName, posx: parseInt(tarx), posy: parseInt(tary), drive: toDrive };
    let reqData = JSON.stringify(movInfo);
    console.log("[moveUser]req:", reqData);
    $.post("/gamectrl/move", reqData, function (data, status) {
        if (status == "success") {
            console.log("moveUser:", data);
            if (data == "ok") {
                // 向所有人广播
                // changePlayerPos(myName, player["posx"], player["posy"], tarx, tary, energy);
            } else {
                alert("无法到达！");
            }
        }
    });
};
var getPlayerObj = function () {
    let playerObj = null;
    $.each(players, function (i, e) {
        // console.log("players e:", e, e["username"]);
        if (e["username"] == myName) {
            playerObj = e;
        }
    });
    return playerObj;
};
// 更新游戏信息
var updateGame = function (req) {
    console.log("updateGame:", req);
    // 更新回合数
    $("#gameround").text(" Round " + req["round"]);
    // 更新当前研究进度
    $("#research").text(req["progress"]);
    // 恢复显示“回合结束”按钮
    $("#btn-finish").show();
    // $("#status").text("");
    // 司机携带人置空
    toDrive = "";
    // 先清空地图上的病毒、玩家和研究所（这里写的不太好）
    for (let i = 0; i < 7; i++) {
        for (let j = 0; j < 7; j++) {
            $("#" + i + "-" + j).children().eq(0).empty();
            $("#" + i + "-" + j).children().eq(1).empty();
        }
    }
    $.each(players, function (i, e) {
        $("#player-" + e["username"]).remove();
    });
    // 更新地图显示
    mapInfo = req["map"];
    for (let i = 0; i < 7; i++) {
        for (let j = 0; j < 7; j++) {
            // 添加病毒
            if (mapInfo[i][j]["virus"] != 0) {
                let virusImg = $('<img>').attr('src', '/static/img/virus.png').attr('class', 'level' + mapInfo[i][j]["virus"]);
                $("#" + i + "-" + j).children().eq(0).append(virusImg);
            }
            // 添加研究所
            if (mapInfo[i][j]["institute"] != 0) {
                let instituteImg = $('<img>').attr('src', '/static/img/cap1.png');
                $("#" + i + "-" + j).children().eq(1).append(instituteImg);
            }
        }
    }
    // 更新玩家位置和职业，放在players全局变量中
    players = req["players"];
    $.each(players, function (i, e) {
        if (e["username"] == myName) {
            myOccupation = occupations[e["occupation"]];
            setCapability(e["occupation"]);
            $("#energy").text(e["energy"]);
        }
        let playerEle = $("<span>").text(e["username"] + "(" + occupations[e["occupation"]] + ")").attr("id", "player-" + e["username"]);
        $("#" + e["posx"] + "-" + e["posy"]).children().eq(2).append(playerEle);
    });
}
// 根据职业设置能力
var setCapability = function(occupationNum) {
    $("#action").text('');
    $("#action").append($("<h2>").text(occupations[occupationNum]));
    $("#action").append($("<p>").text(occupationsDetail[occupationNum]));
    $("#action").append('<a class="toolicon" href="javascript:killVirus();"><img src="/static/img/weapon2.png" /></a>');
    if (occupationNum == 0) { // 科学家
        $("#action").append('<a class="toolicon" href="javascript:doResearch();"><img src="/static/img/weapon1.png" /></a>');
    } else if (occupationNum == 1) { // 工程师
        $("#action").append('<a class="toolicon" href="javascript:buildInstitute();"><img src="/static/img/weapon2.png" /></a>');
    } else if (occupationNum == 2) { // 医生
    } else if (occupationNum == 3) { // 司机
        $("#action").append('<a class="toolicon" href="javascript:drivePlayer();"><img src="/static/img/weapon2.png" /></a>');
    }
}
// 清理病毒
var killVirus = function() {
    console.log("killVirus");
    let player = getPlayerObj();
    if (player == null) {
        return;
    }
    if (mapInfo[player["posx"]][player["posy"]]["virus"] == 0) {
        alert("相同格子内没有病毒");
        return;
    }
    let energy = 1.0;
    if (occupations[player["occupation"]] != "doctor") {
        energy = 1.0 + mapInfo[player["posx"]][player["posy"]]["virus"];
    }
    if (player["energy"] - energy < 0.0) {
        alert("energy is not enough");
        return;
    }
    let killVirusInfo = { username: myName, posx: parseInt(player["posx"]), posy: parseInt(player["posy"]) };
    console.log(killVirusInfo);
    let reqData = JSON.stringify(killVirusInfo);
    $.post("/gamectrl/killvirus", reqData, function (data, status) {
        if (status == "success") {
            console.log("killVirus:", data);
            if (data == "ok") {
                cleanVirus(player["posx"], player["posy"], energy);
            } else {
                alert("清除病毒失败！");
            }
        }
    });
};
var cleanVirus = function(posx, posy, energyNeed) {
    $("#" + posx + "-" + posy).children().eq(0).empty();
    // 修改用户能量
    $.each(players, function (i, e) {
        // console.log("players e:", e, e["username"]);
        if (e["username"] == myName) {
            e["energy"] -= energyNeed;
            // 修改当前能量显示
            $("#energy").text(e["energy"]);
        }
    });
};
// 职业能力
var doResearch = function() {
    console.log("doResearch");
    let player = getPlayerObj();
    if (player == null) {
        return;
    }
    if (occupations[player["occupation"]] != "scientist") {
        return;
    }
    if (player["energy"] - 2.0 < 0.0) {
        alert("energy is not enough");
        return;
    }
    if (mapInfo[player["posx"]][player["posy"]]["institute"] == 0) {
        alert("there is not a institute here");
        return;
    }
    let doResearchInfo = { username: myName, posx: parseInt(player["posx"]), posy: parseInt(player["posy"]) };
    console.log(doResearchInfo);
    let reqData = JSON.stringify(doResearchInfo);
    $.post("/gamectrl/doresearch", reqData, function (data, status) {
        if (status == "success") {
            console.log("doResearch:", data);
            if (data == "ok") {
                addResearch(5, 2.0);
            } else {
                alert("研究疫苗失败！");
            }
        }
    });
};
var addResearch =  function(num, energyNeed) {
    researchProgress += num;
    $("#research").text(researchProgress);
    // 修改用户能量
    $.each(players, function (i, e) {
        // console.log("players e:", e, e["username"]);
        if (e["username"] == myName) {
            e["energy"] -= energyNeed;
            // 修改当前能量显示
            $("#energy").text(e["energy"]);
        }
    });
};
var buildInstitute = function() {
    console.log("buildInstitute");
    let player = getPlayerObj();
    if (player == null) {
        return;
    }
    if (occupations[player["occupation"]] != "engineer") {
        return;
    }
    if (player["energy"] - 4.0 < 0.0) {
        alert("energy is not enough");
        return;
    }
    if (mapInfo[player["posx"]][player["posy"]]["institute"] == 1) {
        alert("there is a institute here");
        return;
    }
    let buildInstituteInfo = { username: myName, posx: parseInt(player["posx"]), posy: parseInt(player["posy"]) };
    console.log(buildInstituteInfo);
    let reqData = JSON.stringify(buildInstituteInfo);
    $.post("/gamectrl/buildinstitute", reqData, function (data, status) {
        if (status == "success") {
            console.log("buildinstitute:", data);
            if (data == "ok") {
                createInstitute(player["posx"], player["posy"], 4.0);
            } else {
                alert("建造研究所失败！");
            }
        }
    });
};
var createInstitute = function(posx, posy, energyNeed) {
    // 修改用户能量
    $.each(players, function (i, e) {
        // console.log("players e:", e, e["username"]);
        if (e["username"] == myName) {
            e["energy"] -= energyNeed;
            // 修改当前能量显示
            $("#energy").text(e["energy"]);
        }
    });
};
var drivePlayer = function() {
    console.log("drivePlayer");
    let player = getPlayerObj();
    if (player == null) {
        return;
    }
    if (occupations[player["occupation"]] != "driver") {
        return;
    }
    if (toDrive != "") {
        toDrive = "";
        alert("放下" + toDrive + "成功！");
        return;
    }
    // todo: 选择一个玩家搭载
    $.each(players, function (i, e) {
        // console.log("players e:", e, e["username"]);
        if ((e["posx"] == player["posx"]) && (e["posy"] == player["posy"]) && (e["username"] != myName)) {
            toDrive = e["username"];
        }
    });
    if (toDrive != "") {
        alert("搭载" + toDrive + "成功！");
        return;
    }
    alert("相同格子内没有其他玩家");
};
// 处理游戏指令
var processCmd = function (req) {
    switch (req["type"]) {
        case cmds.MOV_PLAYER: {
            // 目标格子可能有其他玩家，不能用text()，要用append()
            let playerEle = $("#player-" + req["username"]);
            playerEle.remove();
            $("#" + req["posx"] + "-" + req["posy"]).children().eq(2).append(playerEle);
            // 更新玩家信息
            $.each(players, function (i, e) {
                // console.log("players e:", e, e["username"]);
                if (e["username"] == req["username"]) {
                    e["posx"] = req["posx"];
                    e["posy"] = req["posy"];
                    e["energy"] -= req["energy"];
                    if (req["username"] == myName) {
                        // 修改当前能量显示
                        $("#energy").text(e["energy"]);
                    }
                }
            });
        }break;
        case cmds.SUB_VIRUS: {
            $("#" + req["posx"] + "-" + req["posy"]).children().eq(0).empty();
            mapInfo[req["posx"]][req["posy"]]["virus"] = 0;
        }break;
        case cmds.ADD_INSTITUE: {
            let instituteImg = $('<img>').attr('src', '/static/img/cap1.png');
            $("#" + req["posx"] + "-" + req["posy"]).children().eq(1).append(instituteImg);
            mapInfo[req["posx"]][req["posy"]]["institute"] = 1;
        }break;
        case cmds.DO_RESEARCH: {
            $("#research").text(req["progress"]);
        }break;
        default: {
            console.log("wrong cmd:", req["type"]);
        }
    }
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
        // todo: 回合结束的同步
        let reqData = JSON.parse(e.data);
        if (reqData["users"] !== undefined) { // 收到服务器定时发送的用户信息等
            updateUsers(reqData);
        } else if (reqData["id"] !== undefined) { // 收到消息
            appendMessage(reqData);
        } else if (reqData["running"] !== undefined) { // 开始、结束游戏
            changeGameStatus(reqData);
        } else if (reqData["map"] !== undefined) { // 处理地图信息
            updateGame(reqData);
        } else if (reqData["type"] !== undefined) { // 处理游戏指令
            processCmd(reqData);
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
    $("#btn-ready").on("click", function () {
        if ($(this).text() == "准备") {
            userChangeReady('ready');
        } else {
            userChangeReady('cancel');
        }
    });
    $("#btn-finish").on("click", function() {
        userFinish();
    });
    $(".square").on("click", function () {
        moveUser(this);
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