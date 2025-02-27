function fetchExpression() {
    var expressionId = document.getElementById("expressionId").value;

    if (expressionId === "") {
        alert("Please enter an expression ID.");
        return;
    }

    var xhr = new XMLHttpRequest();
    xhr.open("GET", "http://localhost:8080/api/v1/expressions?id=" + expressionId, true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    var response = JSON.parse(xhr.responseText);
                    displayExpression(response);
                } catch (e) {
                    console.error("Error parsing JSON response:", e);
                }
            } else {
                console.error("Error:", xhr.statusText);
                displayError(xhr.status, xhr.statusText);
            }
        }
    };

    xhr.send();
}

function displayExpression(expression) {
    document.getElementById("button").style.marginBottom = "20px";
    var detailsDiv = document.getElementById("expression-details");
    detailsDiv.innerHTML = ""; // Clear previous content

    var idPara = document.createElement("p");
    idPara.textContent = "ID: " + expression.id;

    var statusPara = document.createElement("p");
    statusPara.textContent = "Status: " + expression.status;
    if (expression.status === 200 || expression.status === 201) {
        statusPara.style.color = "green";
    } else {
        statusPara.style.color = "red";
    }

    var resultPara = document.createElement("p");
    resultPara.textContent = "Result: " + formatResult(expression.result);

    detailsDiv.appendChild(idPara);
    detailsDiv.appendChild(statusPara);
    detailsDiv.appendChild(resultPara);

    // Show the container
    detailsDiv.style.display = "block";
}

function displayError(status, message) {
    document.getElementById("button").style.marginBottom = "20px";
    var detailsDiv = document.getElementById("expression-details");
    detailsDiv.innerHTML = ""; // Clear previous content

    var statusPara = document.createElement("p");
    statusPara.textContent = "Status: " + status;
    statusPara.style.color = "red";

    var resultPara = document.createElement("p");
    resultPara.textContent = "Result: " + message;

    detailsDiv.appendChild(statusPara);
    detailsDiv.appendChild(resultPara);

    // Show the container
    detailsDiv.style.display = "block";
}

function formatResult(result) {
    // Check if the result is a valid number
    var num = parseFloat(result);

    if (isNaN(num)) {
        // If the result is not a number, return it as-is
        return result;
    }

    // Check if the number is an integer
    if (Number.isInteger(num)) {
        return num.toString();
    }

    // Remove trailing zeros
    return num.toFixed(6).replace(/\.?0+$/, '');
}

// Fetch expression when the page loads if ID is provided in the URL
window.onload = function() {
    var urlParams = new URLSearchParams(window.location.search);
    var id = urlParams.get('id');
    if (id) {
        document.getElementById("expressionId").value = id;
        fetchExpression();
    }
};

window.onload = function() {
    document.getElementById("button").style.marginBottom = "0px";
};