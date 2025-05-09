document.getElementById("register-result").style.display = "none"
function register() {
    document.getElementById("register-result").style.display = "block"
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;

    if (username === "" || password === "") {
        alert("Please enter both username and password.");
        return;
    }

    var xhr = new XMLHttpRequest();
    var url = window.location.protocol + "//" + window.location.host + "/api/v1/register";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            
            document.getElementById("register-result").style.display = "block";
            var response = JSON.parse(xhr.responseText);

            if (xhr.status === 200) {
                document.cookie = "token=" + response.token + "; path=/";
                document.getElementById("register-result").textContent = "Registration successful!";
                window.location.replace(window.location.protocol + "//" + window.location.host + "/calculate");
            } else {
                console.error("Error:", response.error);
                document.getElementById("register-result").textContent = "Error: " + response.error;
            }
        }
    };

    var data = JSON.stringify({ login: username, password: password });
    xhr.send(data);
}