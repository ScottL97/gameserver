<!DOCTYPE html>

<html>
<head>
  <title>server</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <script type="text/javascript">
    document.write('<link rel="stylesheet" href="/node_modules/bootstrap/dist/css/bootstrap.css?time=' +
            new Date().getTime() + '"/>');
    document.write('<link rel="stylesheet" href="/static/css/index.css?time=' + new Date().getTime() + '"/>');
    document.write('<link rel="stylesheet" href="/static/css/game.css?time=' + new Date().getTime() + '"/>');
  </script>
</head>

<body>
  <header>
    <h1>server</h1>
    <audio controls="controls" autoplay="autoplay">
      <source src="/static/music/Alpha.mp3" type="audio/ogg" />
      <source src="/static/music/Alpha.mp3" type="audio/mpeg" />
    Your browser does not support the audio element.
    </audio>
    <div>在线用户: <span id="player-num"></span></div>
    <div>游戏状态：<span id="game-status"></span></div>
    <div>已准备：<span id="players"></span></div>
    <div class="you">
      <p>姓名: <span id="username"></span></p>
      <p>行动力: <span id="energy">0</span></p>
      <p>疫苗研究进度：<span id="research">0</span>%</p>
      <!--<p>状态：<span id="status"></span></p>-->
      <div class="form-group">
        <input id="msg" type="text" class="form-control"/>
      </div>
      <div class="text-right">
        <button id="btn-ready" class="btn btn-info">准备</button>
        <button id="btn-finish" class="btn btn-info">回合结束</button>
        <button id="btn-send" class="btn btn-info">发送</button>
      </div>
    </div>
    <div id="your-toolbox" class="toolbox">
      <!--<h2>道具栏</h2>
      <div id="tools">
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon1.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon2.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon3.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon1.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon2.png" /></a>
        <a class="toolicon" href="javascript:;"><img src="/static/img/weapon3.png" /></a>
      </div>-->
      <h2>行动</h2>
      <div id="action"></div>
    </div>
    <div id="your-message" class="list-group"></div>
  </header>
  <div id="gamewindow">
    <div id="gameheader">
      <span id="gamename">游戏名：消毒小队</span>
      <span id="gameround"> Round 1</span>
    </div>
    <div id="game" style="display: none;">
      <div id="0-0" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="0-1" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="0-2" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="0-3" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="0-4" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="0-5" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="0-6" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>

      <div id="1-0" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="1-1" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="1-2" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="1-3" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="1-4" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="1-5" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="1-6" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>

      <div id="2-0" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="2-1" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="2-2" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="2-3" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="2-4" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="2-5" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="2-6" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>

      <div id="3-0" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="3-1" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="3-2" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="3-3" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="3-4" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="3-5" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="3-6" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>

      <div id="4-0" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="4-1" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="4-2" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="4-3" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="4-4" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="4-5" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="4-6" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>

      <div id="5-0" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="5-1" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="5-2" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="5-3" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="5-4" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="5-5" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="5-6" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>

      <div id="6-0" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="6-1" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="6-2" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="6-3" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="6-4" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="6-5" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
      <div id="6-6" class="square">
        <p></p>
        <p></p>
        <p></p>
      </div>
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