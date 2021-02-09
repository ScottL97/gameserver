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
      <div class="form-group">
        <input id="msg" type="text" class="form-control"/>
      </div>
      <div class="text-right">
        <button id="btn-ready" class="btn btn-info">准备</button>
        <button id="btn-send" class="btn btn-info">发送</button>
      </div>
    </div>
    <div id="your-toolbox" class="toolbox">
      <h2>道具栏</h2>
      <div id="tools">
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon1.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon2.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon3.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon1.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon2.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon3.png" /></a>
      </div>
    </div>
    <div id="your-message" class="list-group"></div>
  </header>
  <div id="gamewindow">
    <div id="gameheader">
      <span id="gamename">游戏名：勇士沾恶龙</span>
      <span id="gamecharge">剩余电量：10%</span>
    </div>
  </div>
  <div class="content">
  </div>
  <script src="/node_modules/jquery/dist/jquery.js"></script>
  <script src="/node_modules/popper.js/dist/popper.js" type="module"></script>
  <script src="/node_modules/bootstrap/dist/js/bootstrap.js"></script>
  <script src="/static/js/cookie.js"></script>
  <script src="/static/js/index.js"></script>
</body>
</html>