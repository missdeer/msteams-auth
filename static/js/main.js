function readSettings(accessToken) {
    var xhr = new XMLHttpRequest();
    var url = "https://graph.microsoft.com/v1.0/me/extensions";
    xhr.open("GET", url, false);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Authorization", "Bearer " + accessToken);
    xhr.send();
    if (xhr.readyState === 4 && xhr.status === 200) {
        document.getElementById('settings').value = "got settings from open extensions:\n" + xhr.responseText;
        var jsonVal = JSON.parse(xhr.responseText);
        var values = jsonVal.value;
        var firstVal = values[0];
        var extensionId = firstVal.id
        console.log(firstVal.id)
    } else {
        document.getElementById('settings').value = "read settings response status:" + xhr.status;
    }
}

function deleteExtension(accessToken, extensionId) {
    var xhr = new XMLHttpRequest();
    var url = "https://graph.microsoft.com/v1.0/me/extensions/" + extensionId;
    xhr.open("DELETE", url, false);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Authorization", "Bearer " + accessToken);
    xhr.send();
    if (xhr.readyState != 4 || xhr.status != 204) {
        document.getElementById('settings').value = "delete extension response status:" + xhr.status;
        console.log("delete extension response status:" + xhr.status);
    }
}

function deleteSettings(accessToken) {
    deleteExtension(accessToken, "com.cisco.jabber.integration")
}

function updateExtension(accessToken, extensionId) {
    var xhr = new XMLHttpRequest();
    var url = "https://graph.microsoft.com/v1.0/me/extensions/" + extensionId;
    xhr.open("PATCH", url, false);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Authorization", "Bearer " + accessToken);
    var data = JSON.stringify({
        "@odata.type": "#microsoft.graph.openTypeExtension",
        "id": extensionId,
        "settings": document.getElementById('settings').value
    });
    xhr.send(data);
    if (xhr.readyState != 4 || xhr.status != 204) {
        document.getElementById('settings').value = "delete extension response status:" + xhr.status;
        console.log("delete extension response status:" + xhr.status);
    }
}

function updateSettings(accessToken) {
    updateExtension(accessToken, "com.cisco.jabber.integration")
}

function createSettings(accessToken) {
    //deleteExtension(accessToken, extensionId)
    var xhr = new XMLHttpRequest();
    var url = "https://graph.microsoft.com/v1.0/me/extensions/";
    xhr.open("POST", url, false);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Authorization", "Bearer " + accessToken);
    var data = JSON.stringify({
        "@odata.type": "#microsoft.graph.openTypeExtension",
        "id": "com.cisco.jabber.integration",
        "settings": document.getElementById('settings').value
    });
    xhr.send(data);
    if (xhr.readyState === 4 && xhr.status === 200) {
        document.getElementById('settings').value = "write settings to Microsoft Graph API Open Extensions";
    } else {
        document.getElementById('settings').value = "write settings response status:" + xhr.status;
    }
}