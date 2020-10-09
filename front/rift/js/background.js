chrome.browserAction.onClicked.addListener(function (tab) {
    // var newURL = "index.html";
    // chrome.tabs.create({url: newURL});

    //如果主页已经被打开，则直接跳转过去（不再打开新的页面）
    if (chrome.runtime.openOptionsPage) {
        chrome.runtime.openOptionsPage();
    } else {
        window.open(chrome.runtime.getURL('index.html'));
    }
});


