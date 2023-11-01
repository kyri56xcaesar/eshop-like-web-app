function Login() {
    const getUsernameLogin = document.getElementById("login-username");
    const getPasswordLogin = document.getElementById("login-password");

    console.log(getUsernameLogin.value);
    console.log(getPasswordLogin.value);

    alert(`Recieved: ${getUsernameLogin.value} and ${getPasswordLogin.value}`)
}   

function Register() {
    const getUsernameRegister = document.getElementById("register-username");
    const getPasswordRegister = document.getElementById("register-password");
    const getPasswordRepeatRegister = document.getElementById("register-password-repeat");

    const getEmailRegister = document.getElementById("register-password");
}