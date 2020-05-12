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

function readValue(accessToken) {
    var readReq = new XMLHttpRequest();
    var url = "https://graph.microsoft.com/v1.0/me/extensions/com.cisco.jabber.integration";
    readReq.open("GET", url, false);
    readReq.setRequestHeader("Content-Type", "application/json");
    readReq.setRequestHeader("Authorization", "Bearer " + accessToken);
    readReq.send();
    if (readReq.readyState === 4 && readReq.status === 200) {
        var key = document.getElementById('settings').value
        var jsonVal = JSON.parse(readReq.responseText);
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
    } else {
        throw "Reading Microsoft Graph API Open Extension failed";
    }
}

function writeValue(accessToken) {
    var input = document.getElementById('settings').value
    var inputs = input.split(":")
    if (inputs.length != 2) {
        console.log("invalid inputs")
        return
    }
    var key = inputs[0].trim()
    var value = inputs[1].trim()
    console.log("key: " + key)
    console.log("value: " + value)
    //------------------------------------------------------------------
    var readReq = new XMLHttpRequest();
    var url = "https://graph.microsoft.com/v1.0/me/extensions/com.cisco.jabber.integration";
    readReq.open("GET", url, false);
    readReq.setRequestHeader("Content-Type", "application/json");
    readReq.setRequestHeader("Authorization", "Bearer " + accessToken);
    readReq.send();
    if (readReq.readyState === 4 && readReq.status === 404) {
        // create one
        var createReq = new XMLHttpRequest();
        var url = "https://graph.microsoft.com/v1.0/me/extensions/";
        createReq.open("POST", url, false);
        createReq.setRequestHeader("Content-Type", "application/json");
        createReq.setRequestHeader("Authorization", "Bearer " + accessToken);
        var val = {
            "@odata.type": "#microsoft.graph.openTypeExtension",
            "id": "com.cisco.jabber.integration"
        }
        val[key] = value
        var data = JSON.stringify(val, null, 0);
        createReq.send(data);
        if (createReq.readyState === 4 && (createReq.status === 200 || createReq.status === 201)) {
            console.log("write settings to Microsoft Graph API Open Extensions");
        } else {
            throw new Error("write settings response status:" + createReq.status);
        }
        return
    }
    //------------------------------------------------------------------
    // existed, update it
    var updateReq = new XMLHttpRequest();
    var url = "https://graph.microsoft.com/v1.0/me/extensions/com.cisco.jabber.integration" ;
    updateReq.open("PATCH", url, false);
    updateReq.setRequestHeader("Content-Type", "application/json");
    updateReq.setRequestHeader("Authorization", "Bearer " + accessToken);
    var val = JSON.parse(readReq.responseText)
    val[key] = value
    var data = JSON.stringify(val, null, 0);
    updateReq.send(data);
    if (updateReq.readyState != 4 || updateReq.status != 204) {
        document.getElementById('settings').value = "update extension response status:" + updateReq.status;
        console.log("update extension response status:" + updateReq.status);
    }
}