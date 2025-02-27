function calculate() {
    var expression = document.getElementById("expression").value;

    if (expression === "") {
        alert("Please enter an expression.");
        return;
    }

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:8080/api/v1/calculate", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var response = JSON.parse(xhr.responseText);

                // Remove any existing result paragraphs
                var existingResultPara = document.getElementById("result");
                if (existingResultPara) {
                    existingResultPara.remove();
                }

                // Create and append the result paragraph
                var resultPara = document.createElement("p");
                resultPara.id = "result";
                resultPara.textContent = "ID: " + formatResult(response.id);
                document.querySelector(".container").insertBefore(resultPara, document.getElementById("result-link"));

                // Show the result link
                var resultLinkDiv = document.getElementById("result-link");
                var expressionLink = document.getElementById("expression-link");
                expressionLink.href = "http://localhost:8080/expression?id=" + response.id;
                expressionLink.textContent = "Check the result here!";
                resultLinkDiv.style.display = "block";
            } else {
                console.error("Error:", xhr.statusText);
                // Create and append the error paragraph
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

window.onload = function() {
    document.getElementById("button").style.marginBottom = "0px";
};