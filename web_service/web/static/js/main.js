async function Register(e) {

    e.preventDefault();

    const getUsername = document.getElementById("register-username").value;
    const getEmail = document.getElementById("register-email").value;
    const getPassword = document.getElementById("register-password").value;
    const getRole = document.getElementById("select-role").value;

    try {
        let headersList = {
            "Accept": "*/*",
            "User-Agent": "Thunder Client (https://www.thunderclient.com)",
            "Content-Type": "application/x-www-form-urlencoded"
           }
           
           let bodyContent = "client_secret=HkQiIZCJEC3hSvGoLeTMDGModjWwEfxY&grant_type=client_credentials&client_id=admin-cli";
           
           let response = await fetch("http://localhost:8080/auth/realms/master/protocol/openid-connect/token", { 
             method: "POST",
             body: bodyContent,
             headers: headersList
           });
           

            const adminAccessToken = await response.json();
            console.log(adminAccessToken);     
           
            if(response.ok) {
                const token = adminAccessToken.access_token;

                let headersList2 = {
                    "Accept": "*/*",
                    "User-Agent": "Thunder Client (https://www.thunderclient.com)",
                    "Authorization": "Bearer "+token,
                    "Content-Type": "application/json"
                   }
                   
                   let bodyContent2 = JSON.stringify({
                       "email": getEmail, 
                       "enabled":"true", 
                       "username":getUsername,
                       "attributes": {
                           "client_id": "client-front"
                       },
                       "groups": [getRole],
                       "credentials":[{"type":"password","value":getPassword,"temporary":false}]
                   });
                   
                   let response2 = await fetch("http://localhost:8080/auth/admin/realms/eshop_project/users", { 
                     method: "POST",
                     body: bodyContent2,
                     headers: headersList2
                   });
                   
                   let data = await response2.text();
                   console.log(data);


                   if(response2.ok) {
                        alert('register user is ok');

                        setTimeout(()=> {
                            window.location.href ="http://localhost:5500/login.html";

                        }, 2000);
                   } else {

                   }
                   
            } else {
                console.log("error...")
            } 
           
    } catch (error) {
        console.log(error)
    }


    return false;
}

async function Login(e) {

    e.preventDefault();
    const getUsernameLogin = document.getElementById("login-username").value;
    const getPasswordLogin = document.getElementById("login-password").value;

    try {
        let headersList = {
            "Accept": "*/*",
            "User-Agent": "Thunder Client (https://www.thunderclient.com)",
            "Content-Type": "application/x-www-form-urlencoded"
           }
           
        console.log(getUsernameLogin);
        console.log(getPasswordLogin);
        let bodyContent = "username="+getUsernameLogin+"&password="+getPasswordLogin+"&client_id=client-front&client_secret=FovwjKxMhL4YSjOftIXavlDzkcAqX9Dr&grant_type=password";
        
        let response = await fetch("http://localhost:8080/auth/realms/eshop_project/protocol/openid-connect/token", { 
          method: "POST",
          body: bodyContent,
          headers: headersList
        });
        

        
        if(response.ok) {
            console.log("gsag");
            const login = await response.json();
            const token = login.access_token;

            const decodeToken = await decodeJwt(token);
            console.log(decodeToken);


            // Local Storage or Cookies Create.
            localStorage.setItem("username", decodeToken.preferred_username);
            //localStorage.setItem("password", )
            localStorage.setItem("email", decodeToken.email)
            // localStorage.setItem("role", decodeToken)

        }else {
            const err = await response.json();
            console.log(err);
        }


    } catch (error) {
        console.log(error);
    }

    return false;

}

function decodeJwt(jwtToken) {
    const base64Url = jwtToken.split('.')[1]; // Get the payload part of the JWT
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // Replace Base64 URL encoding characters
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join('')); // Decode Base64 and handle URI component encoding
  
    return JSON.parse(jsonPayload);
  }