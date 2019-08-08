// 页面加载完成执行转盘控件init操作
$(function () {
    DomeWebController.init();
});

// 抽奖 转盘 旋转方式
function LuckyDo(wheel, type) {
    // ajax请求抽奖接口 /lucky
    $.ajax({
        url:"/lucky",
        cache:false,
        // 返回值 json
        dataType:"json",
        // 接口超时时间 1000毫秒 = 1秒
        timeout:1000,
        error:function(request, msg, code) {
            console.log(msg, " , ", code, " , ", request)
            alert("error:" + request + msg + code);
        },
        // ajax请求抽奖接口调用成功
        success:function(data, msg) {
            console.log(data, " , ", msg)
            // 抽奖成功获奖
            if (data.code == 0) {
                // 旋转 旋转到奖品排序位置 1 2 3
                // TODO : NOTICE 后台奖品 displayorder 设置要和转盘旋转情况对应上
                wheel.wheelOfFortune('rotate', data.gift.displayorder, type);
                if (data.gift.gdata != "") {
                    wheel["LuckyMsg"] = "恭喜你中奖了："
                        + data.gift.title
                        + " (" + data.gift.gdata + ")";
                } else {
                    wheel["LuckyMsg"] = "恭喜你中奖了：" + data.gift.title;
                }
            }
            // 没登录
            else if (data.code == 101) {
                alert(data.msg);
                location.href = "/public/index.html";
            }
            // 其他限制条件 超过抽奖次数 ...
            else if (data.code < 200) {
                alert(data.msg);
            }
            // 其他异常 旋转到4幸运奖 旋转类型(指针旋转 转盘旋转 类型不变)
            else {
                wheel.wheelOfFortune('rotate', 4, type);
                // 旋转完毕 赋值 LuckyMsg 给转盘对象新设置该值
                wheel["LuckyMsg"] = "【4】" + data.msg;
            }
        },
        async:true
    });
}

// 显示抽奖结果
function LuckyShow(wheel) {
    alert(wheel["LuckyMsg"])
}

// 定义转盘
DomeWebController = {
    pool: {
        element: {}
    },
    // 调用1个元素
    getEle: function (k) {
        return DomeWebController.pool.element[k];
    },
    // 设置1个元素
    setEle: function (k, v) {
        DomeWebController.pool.element[k] = v;
    },
    // 初始化元素 事件
    init: function () {
        var that = DomeWebController;
        that.inits.element();
        that.inits.event();
        that.build();
    },
    // 设置2个转盘
    inits: {
        element: function () {
            var that = DomeWebController;
            that.setEle("$wheelContainer", $('#wheel_container'));
            that.setEle("$wheelContainer2", $('#wheel_container2'));

        },
        event: function () {
            var that = DomeWebController;

        }
    },
    // 定义2个转盘
    build: function () {
        var that = DomeWebController;
        that.getEle("$wheelContainer").wheelOfFortune({
            'wheelImg': "static/img/wheel_1/wheel.png",//转轮图片
            'pointerImg': "static/img/wheel_1/pointer.png",//指针图片
            'buttonImg': "static/img/wheel_1/button.png",//开始按钮图片
            //'wSide': 400,//转轮边长(默认使用图片宽度)
            //'pSide': 191,//指针边长(默认使用图片宽度)
            //'bSide': 87,//按钮边长(默认使用图片宽度)
            //奖品角度配置{键:[开始角度,结束角度],键:[开始角度,结束角度],......}
            //键是 1等奖 对应角度 2等奖 3等奖 4幸运奖
            'items': {1: [220, 310], 2: [311, 400], 3: [41, 128], 4: [129, 219]},
            //指针图片中的指针角度(x轴正值为0度，顺时针旋转 默认0)
            //角度 270 在1等奖区间
            'pAngle': 270,
            //'type': 'w',//旋转指针还是转盘('p'指针 'w'转盘 默认'p')
            //'fluctuate': 0.5,//停止位置距角度配置中点的偏移波动范围(0-1 默认0.8)
            //'rotateNum': 12,//转多少圈(默认12)
            //'duration': 6666,//转一次的持续时间(默认5000)
            'click': function () {
                // 转盘 转盘旋转
                LuckyDo(that.getEle("$wheelContainer"), 'w')
                // var key = parseInt(Math.random() * 4) + 1;
                // that.getEle("$wheelContainer").wheelOfFortune('rotate', key,'w');
            },
            //点击按钮的回调
            'rotateCallback': function (key) {
                LuckyShow(that.getEle("$wheelContainer"))
                // alert("左:" + key);
            }//转完的回调
        });

        that.getEle("$wheelContainer2").wheelOfFortune({
            'wheelImg': "static/img/wheel_1/wheel.png",//转轮图片
            'pointerImg': "static/img/wheel_1/pointer.png",//指针图片
            'buttonImg': "static/img/wheel_1/button.png",//开始按钮图片
            //'wSide': 400,//转轮边长(默认使用图片宽度)
            //'pSide': 191,//指针边长(默认使用图片宽度)
            //'bSide': 87,//按钮边长(默认使用图片宽度)
            'items': {1: [220, 310], 2: [311, 400], 3: [41, 128], 4: [129, 219]},//奖品角度配置{键:[开始角度,结束角度],键:[开始角度,结束角度],......}
            'pAngle': 270,//指针图片中的指针角度(x轴正值为0度，顺时针旋转 默认0)
            //'type': 'w',//旋转指针还是转盘('p'指针 'w'转盘 默认'p')
            //'fluctuate': 0.5,//停止位置距角度配置中点的偏移波动范围(0-1 默认0.8)
            //'rotateNum': 12,//转多少圈(默认12)
            //'duration': 6666,//转一次的持续时间(默认5000)
            'click': function () {
                // 转盘 指针旋转
                LuckyDo(that.getEle("$wheelContainer2", 'p'))
                // var key = parseInt(Math.random() * 4) + 1;
                // that.getEle("$wheelContainer2").wheelOfFortune('rotate', key, 'p');
            },
            //点击按钮的回调
            'rotateCallback': function (key) {
                LuckyShow(that.getEle("$wheelContainer2"))
                // alert("右:" + key);
            }//转完的回调
        });
    }
};