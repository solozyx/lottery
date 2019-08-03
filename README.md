# lottery

Go抽奖系统


#1 技术依赖
    iris web框架
    redis 客户端
    mysql 客户端
    golang
        https://golang.org/pkg/
    go get github.com/kataras/iris
        https://godoc.org/github.com/kataras/iris#pkg-index
        https://iris-go.com/v10/recipe

    Go结构的标签
        form, json, xorm
        type User struct {
            ID             int64     `json:"id" form:"id"`
            Firstname      string    `json:"firstname" form:"firstname"`
            Username       string    `json:"username" form:"username"`
            HashedPassword []byte    `json:"-" form:"-"`
            CreatedAt      time.Time `json:"created_at" form:"created_at"`
        }
        type PrizeCode struct {
            Id            int    `xorm:"not null pk autoincr comment('虚拟券ID') INT(10)"`
            GiftId        int    `xorm:"not null default 0 comment('奖品ID') index(gift_id) INT(10)"`
            Code          string `xorm:"not null comment('虚拟券内容') index(gift_id) VARCHAR(255)"`
            SysDateline   int    `xorm:"not null default 0 comment('导入时间') INT(10)"`
            SysLastmodify int    `xorm:"not null default 0 comment('最后修改时间') INT(10)"`
            SysStatus     int    `xorm:"not null default 0 comment('状态，0 正常，1 已作废，2 已发放') SMALLINT(5)"`
        }

    Go依赖库
        go get github.com/iris-contrib/httpexpect
            https://godoc.org/github.com/iris-contrib/httpexpect#pkg-index
        go get github.com/go-sql-driver/mysql
        go get github.com/go-xorm/xorm
            http://xorm.io/docs/
            https://github.com/go-xorm/xorm/blob/master/README_CN.md
            http://godoc.org/github.com/go-xorm/xorm
        go get github.com/go-xorm/cmd/xorm
            https://github.com/go-xorm/cmd/blob/master/README_CN.md
            安装后使用 reverse 自动生成代码
            ln -s /usr/local/gopath/bin/xorm /usr/bin/xorm
            cd /usr/local/gopath/src/github.com/go-xorm/cmd/xorm
            xorm reverse mysql "root:***@tcp(10.100.14.21:3306)/activity?charset=utf8" templates/goxorm
            ls -lh ./models/*.go
        go get github.com/gorilla/websocket
        go get gopkg.in/yaml.v2
        go get github.com/gomodule/redigo/redis
        go get git.apache.org/thrift.git/lib/go/thrift

#2 demo程序
    annualMeeting 年会抽奖
    ticket 彩票刮奖
    wechatShake 微信摇一摇
    alipayFu 支付宝五福
    weiboRedPacket 微博抢红包
    wheel 抽奖大转盘

#M.辅助工具
    通过curl打印出来网络请求的各阶段时间
    curl -s -w %{time_namelookup}::%{time_connect}::%{time_starttransfer}::%{time_total}::%{speed_download} "http://www.so.com/"


#安装thrift程序
	brew install thrift
	go get git.apache.org/thrift.git/lib/go/thrift

#生成代码
	cd /private/var/www/go/src/imooc.com/lottery/thrift
	thrift -out .. --gen go lucky.thrift
	thrift -out .. --gen php lucky.thrift

#下载thrift源码，包括各个语言的类库
	http://www.apache.org/dyn/closer.cgi?path=/thrift/0.11.0/thrift-0.11.0.tar.gz
