var domain = localStorage.domain;
$(document).ready(function () {
    $('body').bootstrapMaterialDesign();

    //文件上传
    $("#btn-upload").click(function () {
    });
    //新建文件夹
    $("#btn-new-floader").click(function () {
        parentPath = getParam()
        path = $("#input-floader-name").val()
        if (path == "" || path == undefined) {
            toastr.error("文件夹名不能为空！")
            return
        }

        $.post(domain + "mkdir",
            {
                dir: parentPath + "/" + path,
            },
            function (data, status) {
                if (data.code == 0) {
                    toastr.success(data.response)
                    $('#model-new-floder').modal('hide')
                    fileList(getParam());
                } else {
                    $('#model-new-floder').modal('hide')
                    toastr.error(data.message)
                }
            }
        );
    });
    //新建文件
    $("#btn-new-file").click(function () {
        let parentPath = getParam();
        let fileName = $("#input-new-file-name").val();
        if (fileName === "" || fileName === undefined) {
            toastr.error("文件名不能为空！");
            return
        }
        $.post(domain + "file/new",
            {
                dir: parentPath,
                fileName: fileName
            },
            function (data, status) {
                if (data.code === 0) {
                    toastr.success(data.response);
                    $('#model-new-file').modal('hide');
                    fileList(getParam());
                } else {
                    $('#model-new-file').modal('hide');
                    toastr.error(data.message);
                }
            }
        );
    });
    //离线下载
    $("#btn-down-offline").click(function () {
        $.post(domain + "upload/online",
            {
                'dir': $('#current-path').val(),
                'url': $('#input-file-url').val(),
                'fileName': $('#input-file-name').val()
            },
            function (data, status) {
                if (data.code == 0) {
                    toastr.success(data.response);
                    fileList(getParam());
                    $('#model-off-line-download').modal('hide');
                } else {
                    toastr.error(data.message);
                    $('#model-off-line-download').modal('hide');
                }
            }
        );
    });
    //下载文件
    $("#btn-download").click(function () {
        var checkID = [];//定义一个空数组
        $("input[name='file-id']:checked").each(function (i) {//把所有被选中的复选框的值存入数组
            checkID[i] = $(this).attr("url");
        });
        checkID.forEach(url => {
            const aLabel = document.createElement('a');
            const aEvent = document.createEvent('MouseEvents'); // 创建鼠标事件对象

            aEvent.initEvent('click', false, false);
            aLabel.href = url;
            aLabel.download = ''; // 设置下载文件名，当不需要重新命名时，可以赋值为空字符串
            // 同源文件可以不用写这句，如果非同源文件，请一定加上这句话
            // 否则每次只会下载其中一个文件就中断其余文件的下载了，控制台报和单个文件中截图的那个警告
            // aLabel.target = '_blank';
            aLabel.dispatchEvent(aEvent);
        });
    });
    //删除文件
    $("#btn-delete").click(function () {
        path = getParam();
        var checkID = [];//定义一个空数组
        $("input[name='file-id']:checked").each(function (i) {//把所有被选中的复选框的值存入数组
            checkID[i] = path + "/" + $(this).val();
        });
        $.ajax({
            url: domain + "files",
            type: "delete",
            async: false, // 使用同步方式
            data: JSON.stringify({
                files: checkID
            }),
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            success: function (result) {
                if (result.code == 0) {
                    toastr.success(result.response);
                    $('#model-delete-confirm').modal('hide');
                    fileList(getParam());
                }
            },
            error: function (xhr, status, error) {
                toastr.error(error);
                $('#model-delete-confirm').modal('hide')
            }
        });
    });
    //移动文件
    $("#btn-file-move").click(function () {
        let basePath = getParam();
        let checkID = [];//定义一个空数组
        $("input[name='file-id']:checked").each(function (i) {//把所有被选中的复选框的值存入数组
            checkID[i] = basePath + "/" + $(this).val();
        });
        let movePath = $("#input-file-move-path").val();
        $.ajax({
            url: domain + "move",
            type: "POST",
            async: false, // 使用同步方式
            data: JSON.stringify({
                files: checkID,
                path: [movePath],
            }),
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            success: function (result) {
                if (result.code === 0) {
                    toastr.success(result.response);
                    fileList(basePath);
                    $('#model-file-move').modal('hide');
                } else {
                    toastr.error(result.message);
                    $('#model-file-move').modal('hide');
                }
            },
            error: function (xhr, status, error) {
                toastr.error(error);
                $('#model-file-move').modal('hide');
            }
        });
    });
    //选择要移动的位置
    $("#btn-move-path-select").click(function () {
        selectPath("input-file-move-path", "move-dropdown-menu")
    });
    //重命名
    $("#btn-rename").click(function () {
        console.log($("#input-rename-old-path").val())
        if ($("#input-rename-new-path").val() === undefined || $("#input-rename-new-path").val() === "") {
            toastr.error("文件名不能为空！！！");
            return
        }
        $.post(domain + "rename",
            {
                oldPath: $("#input-rename-old-path").val(),
                newPath: getParam() + "/" + $("#input-rename-new-path").val(),
            },
            function (data, status) {
                if (data.code === 0) {
                    $('#model-rename').modal('hide');
                    toastr.success("保存成功");
                    fileList(getParam());
                } else {
                    toastr.error(data.message)
                }
            }
        );
    });

    //选择上传文件夹
    $("#btn-path-select").click(function () {
        path = $('#input-file-path').prop("value");
        $.ajax({
            url: domain + "list?dir=" + path,
            success: function (result) {
                $("#dropdown-menu").html("");
                $(result.response).each(function (i, data) {
                    if (path == "/") {
                        path = ""
                    }
                    allPath = path + "/" + this.name;
                    if (this.isDir) {
                        pathItem = '<a class="dropdown-item" href="#" name="dropdown-item" value="' + allPath + '"><i class="fa fa-folder-open" aria-hidden="true"></i> ' + this.name + '</a>'
                        $("#dropdown-menu").append(pathItem);
                    }
                });
                //给所有的文件夹绑定点击事件
                $("a[name='dropdown-item']").click(function (e) {
                    //禁用a标签自带的跳转
                    e.preventDefault();
                    $('#input-file-path').prop("value", $(this).attr("value"));

                });
            },
            error: function (xhr, status, error) {
                toastr.error(error);
            }
        });
    });

    //初始化 全选/全不选 按钮
    $("input[name='all-checkbox'][type='checkbox']").click(function () {
        checked = $(this).prop("checked")
        $("input[type='checkbox'][name='file-id']").each(function () {
            $(this).prop('checked', checked);
        })
    });

    $('#btn-init-confirm').click(function () {
        localStorage.domain = $('#input-domain').val();
        localStorage.username = $('#input-username').val();
        localStorage.password = $('#input-password').val();
        localStorage.token = $('#input-token').val();
        //自动加上斜线
        if (localStorage.domain.charAt(localStorage.domain.length - 1) !== "/") {
            localStorage.domain = localStorage.domain + "/"
        }
        domain = localStorage.domain;
        if (domain == undefined || domain == "") {
            toastr.error("domain must not null")
        } else {
            $('#model-init').modal('hide');
            fileList(getParam());
        }
    })

});

if (inited()) {
    fileList(getParam());
}

function inited() {
    toastr.options.closeButton = true;
    if (localStorage.domain === undefined || localStorage.domain == "") {
        $('#model-init').modal('show');
        $('#input-domain').prop("value", localStorage.domain);
        $('#input-username').prop("value", localStorage.username);
        $('#input-password').prop("value", localStorage.password);
        $('#input-token').prop("value", localStorage.token);
        return false
    }
    $.ajaxSetup({
        beforeSend: function (xhr) {
            xhr.setRequestHeader('token', localStorage.token)
        }
    });
    //在接受到数据后做统一处理
    $(document).ajaxSuccess(function (event, request, settings) {
        console.log(request.status);
    });
    $(document).ajaxError(function (event, request, settings) {
        if (request.status / 100 !== 2) {
            // toastr.error("请求失败~！");
        }
    });
    return true
}

function selectPath(inputId, dropMenuId) {
    path = $('#' + inputId).prop("value");
    $.ajax({
        url: domain + "list?dir=" + path,
        success: function (result) {
            $("#" + dropMenuId).html("");
            $(result.response).each(function (i, data) {
                if (path == "/") {
                    path = ""
                }
                allPath = path + "/" + this.name;
                if (this.isDir) {
                    pathItem = '<a class="dropdown-item" href="#" name="dropdown-item" value="' + allPath + '"><i class="fa fa-folder-open" aria-hidden="true"></i> ' + this.name + '</a>';
                    $("#" + dropMenuId).append(pathItem);
                }
            });
            //给所有的文件夹绑定点击事件
            $("a[name='dropdown-item']").click(function (e) {
                //禁用a标签自带的跳转
                e.preventDefault();
                $('#' + inputId).prop("value", $(this).attr("value"));
            });
        },
        error: function (xhr, status, error) {
            toastr.error(error);
        }
    });
}

// initialize with defaults
// $("#input-id").fileinput();
function fileList(path) {
    $('#current-path').prop("value", path)
    $('#input-file-path').prop("value", path)
    generateNav(path)
    $.ajax({
        url: domain + "list?dir=" + path,
        success: function (result) {
            $("#fileList").html("")
            $(result.response).each(function (i, data) {
                if (path == "/") {
                    path = ""
                }
                nameCell = '<a target="_blank" class="text-secondary" href="' + this.url + '"><i class="fa fa-file" aria-hidden="true"></i> ' + this.name + '</a>'
                if (this.isDir) {
                    nameCell = '<a href="#' + path + '/' + this.name + '" name="btn-dir"><i class="fa fa-folder-open" aria-hidden="true"></i> ' + this.name + '</a>'
                }
                trRow = '<tr>\
                <th scope="row">\
                  <div class="checkbox">\
                    <label>\
                      <input type="checkbox" name="file-id" dir="' + this.isDir + '" value="' + this.name + '" url="' + this.url + '">\
                      <span class="checkbox-decorator"><span class="check"></span><div class="ripple-container"></div></span>\
                    </label>\
                  </div>   \
                </th>\
                <td name="name">' + nameCell + '</td>\
                <td>' + this.size + '</td>\
                <td>' + this.time + '</td>\
                <td name="operation" style="width: 10%"></td>\
              </tr>';
                // $("tbody").html(trRow)
                $(trRow).appendTo($("#fileList"))

            });
            //给所有的文件夹绑定点击事件
            $("a[name='btn-dir']").click(function () {
                fileList($(this).attr('href').substring(1));
            });
            //重命名、编辑
            $("#fileList tr").hover(function () {
                $(this).addClass("bg-table-row");
                let checkBox = $(this).children("th").children("span").children("div").children("label").children("input[name='file-id']");
                let isDir = checkBox.attr('dir');
                if (isDir === undefined) {
                    checkBox = $(this).children("th").children("div").children("label").children("input[name='file-id']");
                    isDir = checkBox.attr('dir');
                }
                let path = getParam() + "/" + checkBox.val();
                let name = checkBox.val();
                let butts;
                if (isDir === "true") {
                    butts = '<a name=\'file-rename\' href=\'#\' data-toggle=\'tooltip\' data-placement=\'top\' title=\'重命名\'><li class=\'fa fa-pencil\'></li></a>';
                } else {
                    butts = '<a name=\'file-edit\' href=\'#\' data-toggle=\'tooltip\' data-placement=\'top\' title=\'编辑\'><li class=\'fa fa-file-text\'></li></a>&nbsp;&nbsp;<a name=\'file-rename\' href=\'#\' data-toggle=\'tooltip\' data-placement=\'top\' title=\'重命名\'><li class=\'fa fa-pencil\'></li></a>';
                }
                $(this).children("td[name='operation']").append(butts);
                // $('[data-toggle="tooltip"]').tooltip();
                //编辑
                $("a[name='file-edit']").click(function (e) {
                    //禁用a标签自带的跳转
                    e.preventDefault();
                    chrome.tabs.create({url: "editor.html#" + path});
                });
                //重命名
                $("a[name='file-rename']").click(function (e) {
                    //禁用a标签自带的跳转
                    e.preventDefault();
                    $('#input-rename-old-path').prop("value", path);
                    $('#input-rename-new-path').prop("value", name);
                    $('#model-rename').modal('show');
                });
            }, function () {
                $(this).removeClass("bg-table-row");
                $(this).children("td[name='operation']").html("");
            })
        },
        error: function (xhr, status, error) {
            toastr.error(error);
        }
    });
}

function generateNav(path) {
    $("#breadcrumb").html("");
    $('<li class="breadcrumb-item"><a name="btn-dir" href="#">全部文件</a></li>');
    var pathList = [];
    pathList[0] = {name: "全部文件", path: "/"}
    prePath = ""
    $(path.split("/")).each(function (i, data) {
        if (data == "") {
            return
        }
        prePath = prePath + "/" + data;
        pathList.push({name: data, path: prePath})
    })
    $(pathList).each(function (i, data) {
        $("#breadcrumb").append('<li class="breadcrumb-item"><a name="btn-dir" href="#' + data.path + '">' + data.name + '</a></li>')
    })
}

function getParam() {
    sp = window.location.href.split('#')
    params = ""
    if (sp.length == 2) {
        params = sp[1]
    }
    return decodeURI(params)
}


// with plugin options
$("#input-id").fileinput({
    'showUpload': true,
    'previewFileType': 'any',
    'theme': 'fa',
    'uploadUrl': domain + 'upload',
    // 'mergeAjaxCallbacks': true, //保证ajaxSettings生效
    'ajaxSettings': {
        beforeSend: function (xhr, data) {
            //在请求之前带上token
            xhr.setRequestHeader('token', localStorage.token);
            // data.data.dir = $('#input-file-path').val();
            // console.log(data)
        }
    },
    'uploadExtraData': function () {
        return {
            'token': localStorage.token, // for access control / security
            'dir': $('#input-file-path').val(),
        }
    }
}).on('fileuploaded', function (event, preview, index, fileId) {
    res = preview.response;
    resHtml = '<div class="alert alert-primary" role="alert">\
          <button type="button" class="close" data-dismiss="alert" aria-label="Close">\
            <span aria-hidden="true">&times;</span>\
          </button>\
          <h4 class="alert-heading">' + fileId + ' ' + index + '</h4>\
          <table class="table">\
            <tbody>\
              <tr>\
                <td>url</td>\
                <td>' + res.response + '</td>\
              </tr>\
              <tr>\
                <td>markdown</td>\
                <td>[](' + res.response + ')</td>\
              </tr>\
              </tr>\
            </tbody>\
          </table>\
        </div>';
    $(resHtml).appendTo("#success-text");
}).on('fileuploaderror', function (event, data, msg) {
    // console.log('File Upload Error', 'ID: ' + data.fileId + ', Thumb ID: ' + data.previewId);
}).on('filebatchuploadcomplete', function (event, preview, config, tags, extraData) {
    // console.log('File Batch Uploaded', preview, config, tags, extraData);
});

// with plugin options
$("#input-file-upload").fileinput({
    'showUpload': true,
    'showPreview': false,
    'theme': 'fa',
    'uploadUrl': domain + 'upload',
    // 'mergeAjaxCallbacks': true, //保证ajaxSettings生效
    'ajaxSettings': {
        beforeSend: function (xhr, data) {
            //在请求之前带上token
            xhr.setRequestHeader('token', localStorage.token);
            // data.data.dir = $('#current-path').val();
            // console.log(data)
        }
    },
    'uploadExtraData': function () {
        return {
            'token': localStorage.token, // for access control / security
            'dir': $('#current-path').val(),
        }
    }
}).on('fileuploaded', function (event, preview, index, fileId) {
    res = preview.response;
    toastr.success("success")
}).on('fileuploaderror', function (event, data, msg) {
    toastr.error(msg)
    // console.log('File Upload Error', 'ID: ' + data.fileId + ', Thumb ID: ' + data.previewId);
}).on('filebatchuploadcomplete', function (event, preview, config, tags, extraData) {
    // console.log('File Batch Uploaded', preview, config, tags, extraData);
    $('#model-upload-file').modal('hide');
    fileList(getParam())
});