function fetchExpressions() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", window.location.protocol + "//" + window.location.host + "/api/v1/expressions", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    var response = JSON.parse(xhr.responseText);
                    if (response && Array.isArray(response.expressions)) {
                        displayExpressions(response.expressions);
                    } else {
                        console.error("Invalid response format:", response);
                        displayNoExpressions();
                    }
                } catch (e) {
                    console.error("Error parsing JSON response:", e);
                    displayNoExpressions();
                }
            } else {
                console.error("Error:", xhr.statusText);
                displayNoExpressions();
            }
        }
    };

    xhr.send();
}

function displayExpressions(expressions) {
    var expressionsDiv = document.getElementById("expressions");
    var noExpressionsDiv = document.getElementById("no-expressions");
    var expressionsContainer = document.getElementById("expressions-container");

    expressionsDiv.innerHTML = ""; // Clear previous content

    if (expressions.length === 0) {
        noExpressionsDiv.style.display = "block";
        expressionsContainer.style.display = "none";
        return;
    } else {
        noExpressionsDiv.style.display = "none";
        expressionsContainer.style.display = "block";
    }

    expressions.forEach(function(expression, index) {
        var expressionDiv = document.createElement("div");
        expressionDiv.className = "expression";

        var idStatusDiv = document.createElement("div");
        idStatusDiv.className = "id-status";

        var idPara = document.createElement("p");
        idPara.className = "small-text";
        idPara.textContent = "ID: " + expression.id;

        var statusPara = document.createElement("p");
        statusPara.className = "small-text";
        statusPara.textContent = "Status: " + expression.status;
        if (expression.status === "200" || expression.status === "201") {
            statusPara.style.color = "green";
        } else {
            statusPara.style.color = "red";
        }

        idStatusDiv.appendChild(idPara);
        idStatusDiv.appendChild(statusPara);

        var resultPara = document.createElement("p");
        resultPara.textContent = "Result: " + formatResult(expression.result);

        expressionDiv.appendChild(idStatusDiv);
        expressionDiv.appendChild(resultPara);

        // Remove margin-bottom for the last expression
        if (index === expressions.length - 1) {
            expressionDiv.style.marginBottom = "0";
        }

        expressionsDiv.appendChild(expressionDiv);
    });
}

function displayNoExpressions() {
    var expressionsDiv = document.getElementById("expressions");
    var noExpressionsDiv = document.getElementById("no-expressions");
    var expressionsContainer = document.getElementById("expressions-container");

    expressionsDiv.innerHTML = ""; // Clear previous content
    noExpressionsDiv.style.display = "block";
    expressionsContainer.style.display = "none";
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

// Fetch expressions when the page loads
window.onload = fetchExpressions;