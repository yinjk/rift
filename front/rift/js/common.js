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