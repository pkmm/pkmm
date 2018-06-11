使用Golang做自己想做的
=

使用beego框架
-
* 百度贴吧自动的签到. 
* 验证码的识别.
* 自动登陆教务系统，提取成绩.


部署：使用nginx反向代理

go的服务使用supervisor进行守护

例如： 
```
[program:pkmm]
directory=/root/gopath/src/pkmm
command=/root/gopath/src/pkmm/pkmm
autostart=true
autorestart=true
startsecs=10
user=root
redirect_stderr=true
stdout_logfile=/root/pkmm_creash.log
```

