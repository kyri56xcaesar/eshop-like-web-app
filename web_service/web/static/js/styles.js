
function Login_Style() {

    const getLoginForm = document.getElementById("login-form");
    const getRegisterForm = document.getElementById("register-form");
    const getLoginRegisterDiv = document.getElementById("login-register-form");

    getLoginRegisterDiv.style.display = "block";

    getLoginForm.style.display = "block";
    getRegisterForm.style.display = "none";

    console.log("Login form");


}

function Register_Style() {
    const getLoginForm = document.getElementById("login-form");
    const getRegisterForm = document.getElementById("register-form");
    const getLoginRegisterDiv = document.getElementById("login-register-form");

    getLoginRegisterDiv.style.display = "block";
    getLoginForm.style.display = "none";
    getRegisterForm.style.display = "block";



    const getRegisterUsername = document.getElementById("register-username");
    getRegisterUsername.value = "";

    const getRegisterPassword = document.getElementById("register-password");
    getRegisterPassword.setAttribute("placeholder", "Κωδικός");

    console.log("Register form");

}

function Home_Style() {
    const getLoginForm = document.getElementById("login-form");
    const getRegisterForm = document.getElementById("register-form");
    const getLoginRegisterDiv = document.getElementById("login-register-form");

    getLoginRegisterDiv.style.display = "none";
    getLoginForm.style.display = "none";
    getRegisterForm.style.display = "none";


    // const getLoginRegisterOptions = document.getElementById("login-register-options");
    
    // getLoginRegisterOptions.style.display = "none";

}


 



const sidebar = document.getElementById("sdbar");

function menu_function(x) {


  x.classList.toggle("change")

  sidebar.classList.toggle("close");

  const getLeftHeaders = document.getElementById("left-headers");

  getLeftHeaders.classList.toggle("start");



}