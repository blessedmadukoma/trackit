const api_url = "http://localhost:8000";
// async function getAPI() {
//   const response = await fetch(api_url);
//   const data = await response.json();
//   // console.log(data);

//   // Old index page
//   // const { name, email, age } = data; // JS destructor, to assign the values directly to the specified variables

//   // document.getElementById("name").textContent = name;
//   // document.getElementById("email").textContent = email;
//   // document.getElementById("age").textContent = age;

//   // new index page
//   const { message, status } = data;

//   document.getElementById("status").textContent = status;
//   document.getElementById("message").textContent = message;
// }

// getAPI();

// async function getSignupAPI() {
//   const response = await fetch(api_url + "/auth/signup");
//   const data = await response.json();

//   document.getElementById("signUpStatus").textContent = response.status;
//   document.getElementById("signUpMessage").textContent = data;
// }

// getSignupAPI();

// Send data of the Login Form
const myForm = document.getElementById("myForm");
const displayBtn = document.querySelector("#display");

// Get request
// const loginGETHandler = async () => {
//   try {
//     const response = await fetch(api_url, {
//       method: "GET",
//       headers: {
//         "Content-Type": "Application/json",
//       },
//     });
//     const data = await response.json();
//     if (!response.ok) {
//       console.log("Error: ", data.description);
//       return;
//     }
//     console.log(data);
//   } catch (error) {
//     console.log(error);
//   }
// };

// async function loginPOSTHandler(e) {
//   e.preventDefault();
//   const formD = {
//     username: document.getElementById('username').value,
//     password: document.getElementById('password').value,
//   };
//   console.log(formD);
//   try {
//     const formData = new FormData(myForm);

//     const responseData = await postFormDataAsJson({ api_url, formData });

//     console.log({ responseData });

//     // const response = await fetch(api_url, {
//     //   method: "POST",
//     //   // mode: "no-cors",
//     //   // headers: {
//     //   //   "Content-Type": "application/json",
//     //   // },
//     //   body: formD,
//     //   // body: formData,
//     // });
//     // // console.log(response);
//     // // console.log(response.json(), " : ", response.text());
//     // // const data = await response.text();
//     // // if (!response.ok) {
//     // //   console.log("Error: ", data.description);
//     // //   return;
//     // // }
//     // // console.log(response.json());
//     // // console.log(data);
//   } catch (error) {
//     console.log(error);
//   }
// }


const login = async (e) => {
  e.preventDefault();

  const formData = new FormData(myForm);
  const dataForm = Object.fromEntries(formData.entries());
  // const formD = {
  //   username: document.getElementById("username").value,
  //   password: document.getElementById("password").value,
  // };
  options = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify(dataForm),
    // body: JSON.stringify(formD),
  };
  try {
    const response = await fetch(api_url + "/auth/signin", options);
    const result = await response.json();
    if (!response.ok) {
      console.log("Error: ", result.description);
      return;
    }
    console.log(result);


    document.getElementById("out_status").textContent = response.status;
    document.getElementById("out_fname").textContent = result.firstname;
    document.getElementById("out_lname").textContent = result.lastname;
    document.getElementById("out_email").textContent = result.email;
    document.getElementById("out_mobile").textContent = result.mobile;
    document.getElementById("out_password").textContent = result.password;
  } catch (error) {
    console.log("Error in catch: ", error);
  }
};

const DisplayUsers = async (e) => {
  e.preventDefault();
  try {
    const response = await fetch(api_url + "/user/users");
    const data = await response.json();
    if (!response.ok) {
      console.log("Problem in response");
      return
    }
    console.log(data.length);

    var usersDiv = document.querySelector("#displayUsers");

    // Return to the screen
    for (let index = 0; index < data.length; index++) {
      var div = document.createElement("div");
      div.innerHTML = index+1 + ": " + [data[index].firstname + ", " +data[index].lastname+ ", " +data[index].mobile+ ", " +data[index].email];


      // {
      //   "First Name": data[index].firstname,
      //   "Last Name": data[index].lastname,
      //   Email: data[index].email,
      //   Mobile: data[index].Mobile,
      // }

      usersDiv.appendChild(div);
    }
  } catch (error) {
    console.log(error);
    return
  }
};

myForm.addEventListener("submit", login, false);
// myForm.addEventListener("submit", loginPOSTHandler);
displayBtn.addEventListener("click", DisplayUsers, false);
