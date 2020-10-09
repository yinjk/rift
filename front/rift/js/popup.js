document.querySelector('#go-to-options').addEventListener("click",function() {
	if (chrome.runtime.openOptionsPage) {
	    chrome.runtime.openOptionsPage();
	} else {
	    window.open(chrome.runtime.getURL('index.html'));
	}
});
// $(document).ready(function() { $('body').bootstrapMaterialDesign(); });