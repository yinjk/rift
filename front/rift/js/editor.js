let domain = localStorage.domain;
let changed = false;
$(document).ready(function () {
    toastr.options.closeButton = true;
    let toolbar = [
        'emoji',
        'check',
        'line',
        'insert-before',
        'insert-after',
        'table',
        '|',
        'undo',
        'redo',
        'edit-mode',
        'fullscreen',
        {
            name: 'more',
            toolbar: [
                'outline',
                'preview',
                'code-theme',
                'export',
            ],
        }
    ];
    window.vditor = new Vditor('vditor', {
        _lutePath: "js/vditor/lute.min.js",
        toolbar,
        mode: 'ir',
        height: window.innerHeight + 100,
        outline: true,
        debugger: true,
        typewriterMode: true,
        placeholder: '写你所想!',
        preview: {
            markdown: {
                toc: true,
                mark: true,
            },
        },
        after: function () {
            let str = getParam().split("/");
            let title = str[str.length - 1];
            $("#title").html(title);
            $.get(domain + "text?path=" + getParam(),
                function (data, status) {
                    if (data.code === 0) {
                        vditor.setValue(data.response);
                    } else {
                        toastr.error(data.message)
                    }
                }
            );

        },
        //失去焦点触发
        blur: function (value) {
            if (changed) {
                saveText(value)
            }
        },
        input: function (value, previewElement) {
            changed = true;
            console.log(changed)
        },
        toolbarConfig: {
            pin: true,
        },
        counter: {
            enable: true,
            type: 'text',
        },
        hint: {
            emojiPath: 'http://emojihomepage.com/',
            emojiTail: '<a href="http://emojihomepage.com/" target="_blank">设置常用表情</a>',
        },
        tab: '\t',
    });
    $("#vditor").append(window.vditor);

    //监听 ctrl + s 快捷键
    document.onkeydown = function (e) {
        if ('s' === e.key && (navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey)) {
            e.preventDefault();
            let value = window.vditor.getValue(0);
            saveText(value)
        }
    };

    setInterval("autoSave()", 5000);


    $("#theme-dark").click(function () {
        window.vditor.setTheme('dark', 'dark', 'native');
    });
    $("#theme-light").click(function () {
        window.vditor.setTheme('light', 'light', 'github');
    })

});

function saveText(text) {
    $.post(domain + "text",
        {
            dir: getParam(),
            text: text,
        },
        function (data, status) {
            if (data.code === 0) {
                let now = new Date();
                let hours = now.getHours();
                let minutes = now.getMinutes();
                let seconds = now.getSeconds();
                let nowTime = hours + ":" + minutes + ":" + seconds;
                $("#save-tips").html('<small class="text-success">自动同步于：' + nowTime + '</small>');
            } else {
                $("#save-tips").html('<small class="text-danger">保存失败！</small>');
                toastr.error(data.message)
            }
        }
    );
}

function autoSave() {
    if (changed) {
        let value = window.vditor.getValue(0);
        changed = false;
        saveText(value)
    }
}

function getParam() {
    sp = window.location.href.split('#');
    params = "";
    if (sp.length == 2) {
        params = sp[1]
    }
    return decodeURI(params)
}