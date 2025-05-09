function login() {
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;

    if (username === "" || password === "") {
        alert("Please enter both username and password.");
        return;
    }

    var xhr = new XMLHttpRequest();
    var url = window.location.protocol + "//" + window.location.host + "/api/v1/login";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            document.getElementById("login-result").style.display = "block";
            var response = JSON.parse(xhr.responseText);

            if (xhr.status === 200) {
                document.cookie = "token=" + response.token + "; path=/";
                document.getElementById("login-result").textContent = "Login successful!";               
                window.location.replace(window.location.protocol + "//" + window.location.host + "/calculate");
            } else {    
                console.error("Error:", response.error);
                document.getElementById("login-result").textContent = "Error: " + response.error;
            }
        }
    };

    var data = JSON.stringify({ login: username, password: password });
    xhr.send(data);
}