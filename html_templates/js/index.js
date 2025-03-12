function calculate() {
    var expression = document.getElementById("expression").value;

    if (expression === "") {
        alert("Please enter an expression.");
        return;
    }

    var xhr = new XMLHttpRequest();
    var url = window.location.protocol + "//" + window.location.host + "/api/v1/calculate";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var response = JSON.parse(xhr.responseText);

                var existingResultPara = document.getElementById("result");
                if (existingResultPara) {
                    existingResultPara.remove();
                }

                var resultPara = document.createElement("p");
                resultPara.id = "result";
                resultPara.textContent = "ID: " + formatResult(response.id);
                document.querySelector(".container").insertBefore(resultPara, document.getElementById("result-link"));

                var resultLinkDiv = document.getElementById("result-link");
                var expressionLink = document.getElementById("expression-link");
                expressionLink.href = window.location.protocol + "//" + window.location.host + "/expression?id=" + response.id;
                expressionLink.textContent = "Check the result here!";
                resultLinkDiv.style.display = "block";
            } else {
                console.error("Error:", xhr.statusText);
                var errorPara = document.createElement("p");
                errorPara.id = "result";
                errorPara.textContent = "Error: " + xhr.statusText;
                document.querySelector(".container").appendChild(errorPara);
            }
        }
    };

    var data = JSON.stringify({ expression: expression });
    console.log("Sending data:", data); // Log the data being sent
    xhr.send(data);
}

function fetchExpressions() {
    var xhr = new XMLHttpRequest();
    var url = window.location.protocol + "//" + window.location.host + "/api/v1/expressions";
    xhr.open("GET", url, true);
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
        if (expression.status === 200 || expression.status === 201) {
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

function fetchExpression() {
    var expressionId = document.getElementById("expressionId").value;

    if (expressionId === "") {
        alert("Please enter an expression ID.");
        return;
    }

    var xhr = new XMLHttpRequest();
    var url = window.location.protocol + "//" + window.location.host + "/api/v1/expressions?id=" + expressionId;
    xhr.open("GET", url, true);
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

    detailsDiv.style.display = "block";
}

function displayError(status, message) {
    var detailsDiv = document.getElementById("expression-details");
    detailsDiv.innerHTML = ""; // Clear previous content

    var statusPara = document.createElement("p");
    statusPara.textContent = "Status: " + status;
    statusPara.style.color = "red";

    var resultPara = document.createElement("p");
    resultPara.textContent = "Result: " + message;

    detailsDiv.appendChild(statusPara);
    detailsDiv.appendChild(resultPara);

    detailsDiv.style.display = "block";
}

function formatResult(result) {
    var num = parseFloat(result);

    if (isNaN(num)) {
        // If the result is not a number, return it as-is
        return result;
    }

    if (Number.isInteger(num)) {
        return num.toString();
    }

    return num.toFixed(6).replace(/\.?0+$/, '');
}

window.onload = function() {
    var urlParams = new URLSearchParams(window.location.search);
    var id = urlParams.get('id');
    if (id) {
        document.getElementById("expressionId").value = id;
        fetchExpression();
    }
};