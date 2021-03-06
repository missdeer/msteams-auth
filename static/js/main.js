function readSettings(accessToken) {
    var xhr = new XMLHttpRequest();
    var url = "https://msgraph.yii.li/v1.0/me/extensions";
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
    var url = "https://msgraph.yii.li/v1.0/me/extensions/" + extensionId;
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
    var url = "https://msgraph.yii.li/v1.0/me/extensions/" + extensionId;
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
    var url = "https://msgraph.yii.li/v1.0/me/extensions/";
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

function readValue(accessToken) {
    var xhr = new XMLHttpRequest();
    var url = "https://msgraph.yii.li/v1.0/me/extensions";
    xhr.open("GET", url, false);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Authorization", "Bearer " + accessToken);
    xhr.send();
    if (xhr.readyState === 4 && xhr.status === 200) {
        var key = document.getElementById('settings').value
        var jsonVal = JSON.parse(xhr.responseText);
        var values = jsonVal.value;
        for (var i = 0; i < values.length; i++) {
            var settings = values[i].settings
            if (settings == undefined && settings == null) {
                console.log("find settings field:" + settings);
                return undefined;
            }
            var settingsJSON = JSON.parse(settings);
            var res = settingsJSON[key];
            if (res != undefined && res != null) {
                console.log("find value:" + res);
                return res;
            }
            console.log("can't find value for key: " + key);
            return undefined;
        }
        console.log("can't find value for key: " + key);
        return undefined;
    } else {
        throw "Reading Microsoft Graph API Open Extension failed";
    }
}

function writeValue(accessToken) {
    let data;
    const input = document.getElementById('settings').value;
    const inputs = input.split(":");
    if (inputs.length != 2) {
        console.log("invalid inputs")
        return
    }
    const key = inputs[0].trim();
    const value = inputs[1].trim();
    console.log("key: " + key)
    console.log("value: " + value)
    //------------------------------------------------------------------
    var getReq = new XMLHttpRequest();
    var url = "https://msgraph.yii.li/v1.0/me/extensions";
    getReq.open("GET", url, false);
    getReq.setRequestHeader("Content-Type", "application/json");
    getReq.setRequestHeader("Authorization", "Bearer " + accessToken);
    getReq.send();
    if (getReq.readyState === 4 && getReq.status === 200) {
        var v = hitExtension(getReq.responseText, key)
        if (v != undefined) {
            // found the existed key, modify it
            var settings = v.settings
            var settingsJSON = JSON.parse(settings);
            settingsJSON[key] = value
            var settingsStr = JSON.stringify(settingsJSON, null, 0)
            console.log("settings string: " + settingsStr)

            var updateReq = new XMLHttpRequest();
            var url = "https://msgraph.yii.li/v1.0/me/extensions/" + values[i].id;
            updateReq.open("PATCH", url, false);
            updateReq.setRequestHeader("Content-Type", "application/json");
            updateReq.setRequestHeader("Authorization", "Bearer " + accessToken);
            var data = JSON.stringify({
                "@odata.type": "#microsoft.graph.openTypeExtension",
                "id": values[i].id,
                "settings": settingsStr
            }, null, 0);
            if (data.length > 2048) {
                // try the other extension
            }
            updateReq.send(data);
            if (updateReq.readyState != 4 || updateReq.status != 204) {
                document.getElementById('settings').value = "update extension response status:" + updateReq.status;
                console.log("update extension response status:" + updateReq.status);
            }

            return
        } else {
            throw new Error("write settings response status:" + createReq.status);
        }
        return
    }
    //------------------------------------------------------------------
    // existed, update it
    const updateReq = new XMLHttpRequest();
    url = "https://graph.microsoft.com/v1.0/me/extensions/com.cisco.jabber.integration";
    updateReq.open("PATCH", url, false);
    updateReq.setRequestHeader("Content-Type", "application/json");
    updateReq.setRequestHeader("Authorization", "Bearer " + accessToken);
    const updatedVal = JSON.parse(readReq.responseText);
    if (value == undefined || value == null || value === "") {
        // delete element by key
        delete updatedVal[key];
    } else {
        // update it
        updatedVal[key] = value;
    }
    data = JSON.stringify(updatedVal, null, 0);
    updateReq.send(data);
    if (updateReq.readyState != 4 || updateReq.status != 204) {
        document.getElementById('settings').value = "update extension response status:" + updateReq.status;
        console.log("update extension response status:" + updateReq.status);
    }
}