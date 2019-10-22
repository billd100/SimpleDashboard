document.addEventListener("DOMContentLoaded", function() {
    var button = document.getElementById("div-button").addEventListener("click", addEmptyLinkField);
    
    var inputs = document.getElementsByTagName("input");
    console.log(inputs)
    for (var i = 0; i < inputs.length; i++) {
        inputs[i].addEventListener("keyup", updateNameValueForKeyUp)
    }

    var removes = document.querySelectorAll("td[data-type='remove'");
    for (var i = 0; i < removes.length; i++) {
        removes[i].addEventListener("click", removeElement);
    }
})

function addEmptyLinkField() {
    var inputContainers = document.getElementsByTagName("tr");
    var lastInputContainer = inputContainers[inputContainers.length -1];
    var newInputContainer = lastInputContainer.cloneNode(true);
    var inputs = newInputContainer.getElementsByTagName("input");
    for (var i = 0; i < inputs.length; i++) {
        inputs[i].value = "";
        inputs[i].setAttribute("value", "");
        inputs[i].setAttribute("name", "");
        inputs[i].setAttribute("data-id", "");
        inputs[i].addEventListener("keyup", updateNameValueForKeyUp);
    }
    lastInputContainer.after(newInputContainer);
}

function updateNameValueForKeyUp(e) {
    if (e.target.attributes["data-type"].nodeValue == "name") {
        var inputs = document.querySelectorAll("input[data-id='" + e.target.name + "']");
        console.log(inputs)
        for (var i = 0; i < inputs.length; i++) {
            inputs[i].setAttribute("name", e.target.value);
            inputs[i].setAttribute("data-id", e.target.value);
        }
    } else {
        e.target.setAttribute("value", e.target.value);
    }
}

function removeElement(e) {
    console.log(e.target.parentNode)
    e.target.parentNode.parentNode.removeChild(e.target.parentNode);
}