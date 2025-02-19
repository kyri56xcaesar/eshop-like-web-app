
 function decodeJwt(jwtToken) {
    const base64Url = jwtToken.split('.')[1]; // Get the payload part of the JWT
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // Replace Base64 URL encoding characters
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join('')); // Decode Base64 and handle URI component encoding
  
    return JSON.parse(jsonPayload);
  }

  const client_secr = "RwCn1KFpU2WRDr5DPIX0Ea1usDK0hIER"
  const admin_secr = "SLsbxyA3OhRGE99k0I8lhKkO8UbKGAU5"
  //Login User
  
  async function Login(e){
    //prevent reload page onsubmit
    e.preventDefault()
    //get user username
    const getUsernameLogin = document.getElementById("login-username").value;
    //get user password
    const getPasswordLogin = document.getElementById("login-password").value;
    
    console.log(getUsernameLogin)
    console.log(getPasswordLogin)

    try {
        let headersList = {
            "Accept": "*/*",
            "Content-Type": "application/x-www-form-urlencoded"
           }
           
        let bodyContent = "username="+getUsernameLogin+"&password="+getPasswordLogin+"&client_id=frontend_app&grant_type=password&client_secret="+client_secr;
        
        let response = await fetch("http://localhost:8080/auth/realms/Eshop_project/protocol/openid-connect/token", { 
          method: "POST",
          body: bodyContent,
          headers: headersList
        });
        
        if(response.ok){
            const login = await response.json()
            const token = login.access_token
            const refresh_token = login.refresh_token
            
            //store in localstorage username, email, role (customer, seller) and refresh_token
            const decodeToken = await decodeJwt(token)
            console.log("Decoded token: " + decodeToken.roles)
            console.log(decodeToken.email)
            localStorage.setItem("username", decodeToken.preferred_username)
            localStorage.setItem(decodeToken.preferred_username + "_refreshToken", refresh_token)
            //localStorage.setItem("email", )
            //localStorage.setItem("role", )
            //clear localStorage
            // localStorage.clear()
        }else{
          const err = await response.json()
          console.log(err) 
        }
  
    } catch (error) {
      console.log(error)
    }
    return false
  }
  
  async function Register(e) {
    
    //prevent reload page onsubmit
    e.preventDefault()
  
    //from register form get all data for register a user
    const getUsername = document.getElementById("register-username").value;
    const getEmail = document.getElementById("register-email").value;
    const getPassword = document.getElementById("register-password").value;
    const getRole = document.getElementById("select-role").value;
  
    try {
      var myHeaders = new Headers();
      myHeaders.append("Content-Type", "application/x-www-form-urlencoded");
  
      var urlencoded = new URLSearchParams();
      urlencoded.append("grant_type", "client_credentials");
      urlencoded.append("client_id", "admin-cli");
      urlencoded.append("client_secret", admin_secr);
  
      var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: urlencoded,
        redirect: 'follow'
      };
  
      //get admin access token
      const first_response = await fetch("http://localhost:8080/auth/realms/master/protocol/openid-connect/token", requestOptions)
        
      if(first_response.ok){
        const adminAccessToken = await first_response.json();
        const token = adminAccessToken.access_token
  
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", "Bearer "+token);
  
        var raw = JSON.stringify({
          "email": getEmail,
          "enabled": "true",
          "username": getUsername,
          "attributes": {
            "client_id": "client-front"
          },
          "groups": [
            getRole
          ],
          "credentials": [
            {
              "type": "password",
              "value": getPassword,
              "temporary": false
            }
          ]
        });
  
      var registerOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
      };
  
      const registerUser =  await fetch("http://localhost:8080/auth/admin/realms/Eshop_project/users", registerOptions)
      
      if(registerUser.ok){
        alert('register user is ok')
  
        setTimeout(()=>{
          window.location.href = "http://localhost:8079"
        },2000)
        
      }else{
        const err = await registerUser.json()
        console.log(err)
      }
  
      }else{
        const err = await first_response.json();
        console.log(err);
      }
      
    } catch (error) {
      console.log(error)
    }
  
    return false
  }
  
  
async function Logout(e) {

}
  