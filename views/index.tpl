<!DOCTYPE html>

<html>
<head>
  <title>server</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <script type="text/javascript">
    document.write('<link rel="stylesheet" href="/node_modules/bootstrap/dist/css/bootstrap.css?time=' +
            new Date().getTime() + '"/>');
    document.write('<link rel="stylesheet" href="/static/css/index.css?time=' + new Date().getTime() + '"/>');
  </script>
</head>

<body>
  <header>
    <h1>server</h1>
    <div>在线用户: <span id="player-num"></span></div>
    <div class="you">
      <p>姓名: <span id="username"></span></p>
      <form>
        <div class="form-group">
          <input id="msg" type="text" class="form-control"/>
        </div>
      </form>
      <div class="text-right">
        <button id="btn-ready" class="btn btn-info">准备</button>
        <button id="btn-send" class="btn btn-info">发送</button>
      </div>
    </div>
    <ul id="your-message"></ul>
  </header>
  <div id="gamewindow">
    game on!
  </div>
  <div class="content">
  </div>
  <script src="/node_modules/jquery/dist/jquery.js"></script>
  <script src="/node_modules/popper.js/dist/popper.js" type="module"></script>
  <script src="/node_modules/bootstrap/dist/js/bootstrap.js"></script>
  <script src="/static/js/cookie.js"></script>
  <script>
    var wsProcess = function(myId, myUsername) {
      $("#username").text(myUsername);
      let ws = new WebSocket('ws://' + window.location.host + '/ws');
      let ids = new Map();
      // 建立连接后发送空消息带id，建立用户和连接的关系，断开连接时清除登录状态
      ws.onopen = function(e) {
        console.log("websocket opened!");
        let messageInfo = {id: myId, msg: "", username: myUsername};
        ws.send(JSON.stringify(messageInfo));
      }
      // 服务器WebSocket消息处理
      ws.onmessage = function(e) {
        // console.log(e.data);
        // todo: e.data类型判断
        let message = JSON.parse(e.data);
        if (message["users"] !== undefined) { // 收到服务器定时发送的用户状态信息等
          let users = message["users"];
          // console.log(users);
          let userNum = Object.keys(users).length;
          $("#player-num").text(userNum);
          // 删除不在线的用户
          $(".content").children().each(function(i, e) {
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
              let toolBoxDiv = $('<div>').text("道具栏：");
              userBoxDiv.append(toolBoxDiv);
              $(".content").append(userBoxDiv);
            }
          }
        } else { // 收到消息
          if (message["id"] !== myId) {
            let eleId = "#" + message["id"];
            console.log("message from", eleId);
            $('<p>').text(message["msg"]).appendTo($(eleId));
            if ($(eleId).children().length > 5) {
              $(eleId).children().eq(2).remove();
            }
          }
        }
      };
      // 每次发消息时在服务器校验id
      $("#btn-send").click(function() {
        let msg = $("#msg").val();
        if (msg.length > 0) {
          let messageInfo = {id: myId, msg: msg, username: myUsername};
          ws.send(JSON.stringify(messageInfo));
        } else {
          alert("消息不能为空！");
        }
      });
    };
    $(function() {
      let username = getCookie("username");
      let id = getCookie("token");
      console.log("id:", id);
      console.log("username:", username);
      let userInfo = {username: username, id: id};
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
    
  </script>
</body>
</html>