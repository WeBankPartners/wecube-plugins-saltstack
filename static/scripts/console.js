/**
 *
 * Auth:Eric Shi
 * Mail:postmaster@apibox.club
 * QQ:155122504
 *
 */
(function (factory) {
    "use strict";
    if (typeof define === "function" && (define.amd || define.cmd)) {
        define(["jquery"], factory);
    } else {
        factory((typeof (jQuery) != "undefined") ? jQuery : window.Zepto);
    }
}
    (function ($) {
        "use strict";
        $.extend({
            LogOut: function () {
                location.href = "/console/logout" + "?rnd=" + Math.random();
            }
        });

        $.fn.ConsoleLogin = function (options) {
            if (options === undefined) {
                options = {};
            }
            var loginBoxMsg, username, userpwd, standalone, src_vmaddr, dst_vmaddr, $loginForm = this;
            loginBoxMsg = options.loginBoxMsg === undefined ? $("#loginBoxMsg") : $(options.loginBoxMsg);
            username = options.username === undefined ? $("#username") : $(options.username);
            userpwd = options.userpwd === undefined ? $("#userpwd") : $(options.userpwd);
            src_vmaddr = options.src_vmaddr === undefined ? $("#src_vmaddr") : $(options.src_vmaddr);
            dst_vmaddr = options.dst_vmaddr === undefined ? $("#dst_vmaddr") : $(options.dst_vmaddr);
            standalone = options.standalone === undefined ? false : options.standalone;

            username.focus();

            var boxMsg = function (t) {
                loginBoxMsg.html(t);
                loginBoxMsg.show();
            };
            var p = null;
            $loginForm.bind("form-pre-serialize", function (event, form, options, veto) {
                p = userpwd.val();
            });

            if (standalone === true) {
                var keyEvt = function (obj) {
                    $.post("/console/chksshdaddr?rnd=" + Math.random(), {
                        "vm_addr": obj.val()
                    }, function (data) {
                        console.log("data:", data);
                        var json = data;
                        if (typeof (data) != "object") {
                            json = $.parseJSON(data);
                        }
                        if (json.ok) {
                            console.log(json.data.sshd_addr);
                            dst_vmaddr.val(json.data.sshd_addr);
                        }
                    });
                };
                src_vmaddr.unbind("keyup").keyup(function (evt) {
                    keyEvt($(this));
                });
                src_vmaddr.unbind("paste").bind("paste", function (evt) {
                    keyEvt($(this));
                });
                src_vmaddr.unbind("blur").blur(function (evt) {
                    keyEvt($(this));
                });
            }

            $loginForm.ajaxForm({
                dataType: "json",
                beforeSubmit: function (a, f, o) {
                    loginBoxMsg.hide();
                    if (standalone === true) {
                        if (src_vmaddr.val().length === 0) {
                            src_vmaddr.focus();
                            boxMsg("请输入主机地址");
                            return false;
                        }

                        var u = url.parse(src_vmaddr.val());
                        if (u.host === undefined || u.port === undefined) {
                            src_vmaddr.focus();
                            boxMsg("请输入正确的主机地址");
                            return false;
                        }
                    }
                    if (username.val().length === 0) {
                        userpwd.val(p);
                        username.focus();
                        boxMsg("请输入您用户名");
                        return false;
                    }
                    if (userpwd.val().length === 0) {
                        userpwd.val(p);
                        username.focus();
                        boxMsg("请输入您的密码");
                        return false;
                    }
                },
                success: function (data) {
                    var json = data;
                    if (typeof (data) != "object") {
                        json = $.parseJSON(data);
                    }
                    if (json.ok) {
                        location.href = json.data;
                    } else {
                        userpwd.focus();
                        userpwd.select();
                        boxMsg(json.msg);
                    }
                }
            });
        };

        $.fn.OpenTerminal = function (options) {
            if (options === undefined) {
                options = {};
            }
            var wsaddr, $console = this;
            wsaddr = options.wsaddr === undefined ? "ws://127.0.0.1:8080/console" : options.wsaddr;

            var resizeTerminal = function (t, c, r) {
                var appbar_height = $("#c-appbar").height();
                var body_height = $(window).height();
                var body_width = $(window).width();
                var terminal_height = body_height - appbar_height - 30;
                $(".terminal").height(terminal_height);
                t.resize(c, r);
            };

            var getSize = function () {
                function getCharSize() {
                    var span = $("<span>", { text: "qwertyuiopasdfghjklzxcvbnm" });
                    $console.append(span);
                    var lh = span.css("line-height");
                    lh = lh.substr(0, lh.length - 2);
                    var size = {
                        width: span.width() / 22,
                        height: span.height() - (lh / 2.8)
                    };
                    span.remove();
                    return size;
                }

                function getwindowSize() {
                    var appbar_height = $("#c-appbar").height();
                    var body_height = $(window).height();
                    var body_width = $(window).width() - 20;
                    var terminal_height = body_height - appbar_height -30;
                    return {
                        width: body_width,
                        height: terminal_height
                    };
                }
                var charSize = getCharSize();
                var windowSize = getwindowSize();

                return {
                    cols: Math.floor(windowSize.width / charSize.width),
                    rows: Math.floor(windowSize.height / charSize.height) - 5
                };
            };

            window.WebSocket = window.WebSocket || window.MozWebSocket;
            var cols = getSize().cols;
            var rows = getSize().rows;
            var term = null;
	    console.log(wsaddr)
	    console.log(wsaddr + "&cols=" + cols + "&rows=" + rows)
            var socket = new WebSocket(wsaddr + "&cols=" + cols + "&rows=" + rows);

            Terminal.applyAddon(fit);
            Terminal.applyAddon(attach);

            socket.onopen = function () {
                term = new Terminal({
                    termName: "xterm",
                    cols: cols,
                    rows: rows,
                    useStyle: true,
                    convertEol: true,
                    screenKeys: true,
                    cursorBlink: false,
                    visualBell: true,
                    colors: Terminal.xtermColors
                });

                console.log(term);

                term.attach(socket);
                term._initialized = true;

                term.open($console.get(0));
                term.fit();

                resizeTerminal(term, cols, rows);

                $(window).resize(function () {
                    resizeTerminal(term, cols, rows);
                });

                term.on("title", function (title) {
                    $(document).prop("title", title);
                });

                window.term = term;
                window.socket = socket;
            };
            socket.onclose = function (e) {
                term.destroy();
                var span = $("<span></span>");
                span.text("网络中断（与服务器端连接已断开，请重新连接或联系管理员）。");
                span.css({
                    "display": "inline-block",
                    "position": "absolute",
                    "z-index": 1500,
                    "top": "40%",
                    "width": "100%",
                    "text-align": "center",
                    "color": "#FFFFFF"
                });
                span.appendTo($console);
            };
            socket.onerror = function (e) {
                console.log("Socket error:", e);
            };
        };

        $.fn.loading = function (options) {
            if (options === undefined) {
                options = {};
            }
            var text, ico_class, $loadDiv = this;

            text = options.text === undefined ? "Loading..." : options.text;
            ico_class = options.ico_class === undefined ? "mif-spinner4 mif-3x mif-ani-spin fg-green" : options.ico_class;
            var overlay = $('<div><i class="' + ico_class + '"></i> ' + text + '</div>');
            $loadDiv.html("");
            $loadDiv.css({
                "position": "relative"
            });
            $loadDiv.append(overlay);
            $(overlay).css({
                "display": "inline-block",
                "position": "absolute",
                "z-index": 1500,
                "top": "40%",
                "left": "45%"
            });
        };
    }));
