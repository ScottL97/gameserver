<!DOCTYPE html>

<html>
<head>
  <title>server</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <script type="text/javascript">
    document.write('<link rel="stylesheet" href="/node_modules/bootstrap/dist/css/bootstrap.css?time=' +
            new Date().getTime() + '"/>');
    document.write('<link rel="stylesheet" href="/static/css/login.css?time=' + new Date().getTime() + '"/>');
  </script>
</head>

<body>
  <div class="login-box">
    <div class="form-group">
      <label for="username">User Name</label>
      <input type="text" class="form-control" id="username" aria-describedby="usernameHelp">
      <small id="usernameHelp" class="form-text text-muted">Please input your user name.</small>
    </div>
    <div class="form-group">
      <label for="password">Password</label>
      <input type="password" class="form-control" id="password">
    </div>
    <div class="text-right">
      <button id="loginBtn" class="btn btn-primary">Login</button>
    </div>
  </div>

  <script src="/node_modules/jquery/dist/jquery.js"></script>
  <script src="/node_modules/popper.js/dist/popper.js" type="module"></script>
  <script src="/node_modules/bootstrap/dist/js/bootstrap.js"></script>
  <script src="/static/js/cookie.js"></script>
  <script type="text/javascript">
    $(function() {
      $("#loginBtn").click(function() {
        let usernameInput = $("#username").val();
        let passwordInput = $("#password").val();
        let userInfo = {username: usernameInput, password: passwordInput};
        // 用户登录
        $.post('/user', JSON.stringify(userInfo), function (data, status) {
          if (status == "success") {
            let id = data;
            if (id != "") {
              setCookie("username", usernameInput, 1);
              setCookie("id", id, 1);
              $(window).attr("location", "/game");
            }
          }
        });
      });
    });
  </script>
</body>