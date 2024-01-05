const client_secr = "RwCn1KFpU2WRDr5DPIX0Ea1usDK0hIER"


const user = localStorage.getItem("username");

console.log("Current user: " + user);


if (user == "" || user == null || user == undefined) {
  window.location.href = "http://localhost:5500/web_service/web/error.html";

} else {
  document.getElementById("user").innerHTML = user;
}

async function Logout(e, user) {
  console.log("Logging out "+user+" ...");


  const refresh_token = localStorage.getItem("refreshToken");

  try {
    let headersList = {
        "Accept": "*/*",
        "Content-Type": "application/x-www-form-urlencoded"
       }
       
    // console.log(refresh_token);
    let bodyContent = "refresh_token="+refresh_token+"&client_id=frontend_app&client_secret="+client_secr;
    
    let response = await fetch("http://localhost:8080/auth/realms/Eshop_project/protocol/openid-connect/logout", { 
      method: "POST",
      body: bodyContent,
      headers: headersList
    });

    // console.log(response);
    
    if(response.ok){

        //clear localStorage
        localStorage.clear()

        setTimeout(()=>{
          window.location.href = "http://localhost:5500/web_service/web/"
        }, 500);


    }else{
      const err = await response.json()
      console.log(err) 
    }

} catch (error) {
  console.log(error)
}
return false

}


async function Profile(e) {
  console.log("Profile.");
}

