const client_secr = "HBif856y3ZZcB8NN2KQXqLZJkgu70D7L";
const admin_secr = "kJ3UhrtmNrURm8rX0cjS4fMTM7gXfLug";
const host='localhost';


const l_username = localStorage.getItem("username");
const l_role = localStorage.getItem("role");
console.log("Current user: " + l_username);
console.log("Current role: " + l_role)

if (l_username != "" && l_username != null && l_username != undefined && l_role != "" && l_role != null && l_role != undefined) {
  window.location.href = "http://"+host+":8079/"+l_role+"/";

} 
 
  
 
 
 function decodeJwt(jwtToken) {
    const base64Url = jwtToken.split('.')[1]; // Get the payload part of the JWT
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // Replace Base64 URL encoding characters
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join('')); // Decode Base64 and handle URI component encoding
  
    return JSON.parse(jsonPayload);
  }

  


  var role = ""
  //Login User
  
  async function Login(e){
    //prevent reload page onsubmit
    e.preventDefault()
    //get user username
    const getUsernameLogin = document.getElementById("login-username").value;
    //get user password
    const getPasswordLogin = document.getElementById("login-password").value;
    
    // console.log(getUsernameLogin)
    // console.log(getPasswordLogin)

    role = localStorage.getItem("role");

    if (localStorage.getItem("username") == getUsernameLogin && role != "") {
      console.log("User already logged in.");
      window.location.href = "http://"+host+":8079/"+role+"/"

      return;
    }

    try {

        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/x-www-form-urlencoded");
        myHeaders.append("Accept", "*/*");

        var urlencoded = new URLSearchParams();
        urlencoded.append("username", getUsernameLogin);
        urlencoded.append("password", getPasswordLogin);
        urlencoded.append("client_id", "frontend_app");
        urlencoded.append("grant_type", "password");
        urlencoded.append("client_secret", client_secr);

        var requestOptions = {
          method: 'POST',
          headers: myHeaders,
          body: urlencoded,
          redirect: 'follow'
        };

        let response = await fetch("http://"+host+":8080/auth/realms/Eshop_project/protocol/openid-connect/token", requestOptions);
        
        

        // console.log(response);
        
        if(response.ok){
            const login = await response.json()
            // console.log(login)
            const token = login.access_token
            const refresh_token = login.refresh_token
            
            //store in localstorage username, email, role (customer, seller) and refresh_token
            const decodeToken = await decodeJwt(token)

            localStorage.setItem("username", decodeToken.preferred_username)
            localStorage.setItem("refreshToken", refresh_token)
            localStorage.setItem("email", decodeToken.email)
            

            
            if (decodeToken.realm_access.roles.includes("seller")) {
                console.log("User is a seller");
                role = "seller";


            }
            else if (decodeToken.realm_access.roles.includes("customer")) {
                console.log("User is a customer");
                role = "customer";

            }

            localStorage.setItem("role", role)


            setTimeout(()=>{
              window.location.href = "http://"+host+":8079/"+role+"/"
            }, 500);
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
      myHeaders.append("Accept", "*/*");
  
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
      const first_response = await fetch("http://"+host+":8080/auth/realms/master/protocol/openid-connect/token", requestOptions)
        
      if(first_response.ok){
        const adminAccessToken = await first_response.json();
        const token = adminAccessToken.access_token
  
        var myHeaders = new Headers();
        myHeaders.append("Accept", "*/*");
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
  
      const registerUser =  await fetch("http://"+host+":8080/auth/admin/realms/Eshop_project/users", registerOptions)
      
      if(registerUser.ok){
  

        Login_Style();
        
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
  
  

  