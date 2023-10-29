const getLoginSelect = document.getElementById("login-select");
const getRegisterSelect = document.getElementById("register-select");

const getLoginPassword = document.getElementById("login-password");
const getRegisterUsername = document.getElementById("register-username");
const getRegisterPassword = document.getElementById("register-password");
const getRegisterEmail = document.getElementById("register-email");



function Login_Style() {

    const getLoginForm = document.getElementById("login-form");
    const getRegisterForm = document.getElementById("register-form");
    const getLoginRegisterDiv = document.getElementById("login-register-form");

    getLoginRegisterDiv.style.display = "block";

    getLoginForm.style.display = "block";
    getRegisterForm.style.display = "none";




}

function Register_Style() {
    const getLoginForm = document.getElementById("login-form");
    const getRegisterForm = document.getElementById("register-form");
    const getLoginRegisterDiv = document.getElementById("login-register-form");

    getLoginRegisterDiv.style.display = "block";
    getLoginForm.style.display = "none";
    getRegisterForm.style.display = "block";


    // resetForm(getRegisterForm);


}

function Home_Style() {
    const getLoginForm = document.getElementById("login-form");
    const getRegisterForm = document.getElementById("register-form");
    const getLoginRegisterDiv = document.getElementById("login-register-form");

    getLoginRegisterDiv.style.display = "none";
    getLoginForm.style.display = "none";
    getRegisterForm.style.display = "none";



}

