
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


    const getRegisterUsername = document.getElementById("register-username");
    getRegisterUsername.value = "";

    const getRegisterPassword = document.getElementById("register-password");
    getRegisterPassword.setAttribute("placeholder", "Κωδικός");


    // resetForm(getRegisterForm);


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


 
function menu_function(x) {
    x.classList.toggle("change")

    sidebar.classList.toggle("close");

    const getLeftHeaders = document.getElementById("left-headers");

    getLeftHeaders.classList.toggle("start");

}



const sidebar = document.querySelector(".sidebar");
const menu = document.querySelector(".menu-content");
const menuItems = document.querySelectorAll(".submenu-item");
const subMenuTitles = document.querySelectorAll(".submenu .menu-title");


menuItems.forEach((item, index) => {
  item.addEventListener("click", () => {
    menu.classList.add("submenu-active");
    item.classList.add("show-submenu");
    menuItems.forEach((item2, index2) => {
      if (index !== index2) {
        item2.classList.remove("show-submenu");
      }
    });
  });
});

subMenuTitles.forEach((title) => {
  title.addEventListener("click", () => {
    menu.classList.remove("submenu-active");
  });
});