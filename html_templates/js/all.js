function log_out() {
    //document.cookie.split(";").forEach(function(c) { if (c.split("=")[0] === "token") {document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/"); }});
    delete_cookie("token")
    window.location.pathname = "/login"
}

function set_cookie(name, value) {
    document.cookie = name +'='+ value +'; Path=/;';
}
function delete_cookie(name) {
    document.cookie = name +'=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}