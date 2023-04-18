package mnms

import (
	"reflect"
	"testing"
)

// TestMNMSConfigReadWrite test mnms config read write
func TestMNMSConfigReadWrite(t *testing.T) {
	var err error
	QC.OwnPublicKeys, err = GenerateOwnPublickey()
	if err != nil {
		t.Fatal("generate own public key fail", err)
	}
	// clear mnms config
	err = cleanMNMSConfig()
	if err != nil {
		t.Log("clean mnms config fail", err)
	}

	// Init default mnms config
	err = InitDefaultMNMSConfigIfNotExist()
	if err != nil {
		t.Fatal("init default mnms config fail", err)
	}

	expectConfig := &MNMSConfig{
		Users: []UserConfig{
			{
				Name:     "admin",
				Role:     MNMSAdminRole,
				Password: AdminDefaultPassword,
			},
		},
	}
	// Test readconfig
	readedConfig, err := GetMNMSConfig()
	if err != nil {
		t.Fatal("read mnms config fail", err)
	}
	if !reflect.DeepEqual(readedConfig, expectConfig) {
		t.Fatal("expect mnms config", expectConfig, "but got", readedConfig)
	}

	// modify config
	expectConfig.Users = append(expectConfig.Users, UserConfig{
		Name:     "test",
		Role:     MNMSSuperUserRole,
		Password: "testrawpass",
	})

}

// TestCheckUsersPassword
func TestCheckUsersPassword(t *testing.T) {
	password := "MyPassword123$"
	err := checkUsersPassword(password)
	if err != nil {
		t.Fatal("checkUsersPassword fail", err)
	}
	tooShort := "abcd"
	err = checkUsersPassword(tooShort)
	if err == nil {
		t.Fatal("checkUsersPassword should fail")
	}
	nodigit := "MyPassword$"
	err = checkUsersPassword(nodigit)
	if err == nil {
		t.Fatal("checkUsersPassword should fail")
	}
	nospecial := "MyPassword123"
	err = checkUsersPassword(nospecial)
	if err == nil {
		t.Fatal("checkUsersPassword should fail")
	}
	noUpper := "mypassword123$"
	err = checkUsersPassword(noUpper)
	if err == nil {
		t.Fatal("checkUsersPassword should fail")
	}
}

// TestDeleteUser test delete user
func TestDeleteUser(t *testing.T) {

	// Init default mnms config
	err := InitDefaultMNMSConfigIfNotExist()
	if err != nil {
		t.Fatal("init default mnms config fail", err)
	}

	// delete user
	err = DeleteUserConfig("shouldnotexist")
	if err == nil {
		t.Fatal("delete user should fail")
	}
	// add test user
	err = AddUserConfig("test", MNMSSuperUserRole, "tA%@18632Nest", "test")
	if err != nil {
		t.Fatal("add user fail", err)
	}
	c, err := GetUserConfig("test")
	if err != nil {
		t.Fatal("get user fail", err)
	}
	if c.Name != "test" {
		t.Fatal("get user fail", err)
	}
	// delete user
	err = DeleteUserConfig("test")
	if err != nil {
		t.Fatal("delete user fail", err)
	}
	_, err = GetUserConfig("test")
	if err == nil {
		t.Fatal("get user should fail")
	}
	conf, err := GetMNMSConfig()
	if err != nil {
		t.Fatal("get mnms config fail", err)
	}
	t.Log(conf)
}

// TestDuplicateUserAndEmail test duplicate user and email
func TestDuplicateUserAndEmail(t *testing.T) {
	// Init default mnms config
	err := InitDefaultMNMSConfigIfNotExist()
	if err != nil {
		t.Fatal("init default mnms config fail", err)
	}

	// add test user
	err = AddUserConfig("test", MNMSSuperUserRole, "tA%@18632Nest", "test@test.com")
	if err != nil {
		t.Fatal("add user fail", err)
	}
	// add duplicate user
	err = AddUserConfig("test", MNMSSuperUserRole, "tA%@18632Nest", "other@test.com")
	if err == nil {
		t.Fatal("add duplicate user should fail")
	}
	// add duplicate email
	err = AddUserConfig("test2", MNMSSuperUserRole, "tA%@18632Nest", "test@test.com")
	if err == nil {
		t.Fatal("add duplicate email should fail")
	}
	// delete user
	err = DeleteUserConfig("test")
	if err != nil {
		t.Fatal("delete user fail", err)
	}
}

/*
# Two-Factor Authentication (2FA) Implementation Documentation

## Introduction

We have successfully implemented Two-Factor Authentication (2FA) in our system to enhance the security of user accounts. This document outlines the changes and new API endpoints related to this feature.

## Changes to the User Model

1. Two new fields have been added to the User model: `enable2FA` (boolean) and `secret` (string).
2. By default, `enable2FA` is set to `false`, and `secret` is set to an empty string (`""`).

## API Endpoints and Behavior

1. **`/login`**:
   - If a user logs in successfully and has not enabled 2FA, the API returns a JWT token.
   - If the user has enabled 2FA, the API returns a `sessionID` and the user account instead of a token.
New response body:
```json
{
    "sessionID": "sessionID",
    "user": "user name",
    "email":"user's email"
}
```

2. **`/2fa/validate`**:
   - After a successful `/login` request, users have five minutes to validate their 2FA code through this endpoint.
   - If the validation is successful, the API returns a JWT token.

3. **`POST /2fa/secret`**:
   - To enable 2FA for a user, make a POST request to this endpoint with the following JSON body: `{"user":"user name"}`.
   - The API returns the secret and enables 2FA for the specified user account.

4. **`PUT /2fa/secret`**:
   - To update the 2FA secret for a user, make a PUT request to this endpoint with the following JSON body: `{"user":"user name"}`.
   - The API updates the secret for the specified user account.

5. **`DELETE /2fa/secret`**:
   - To disable 2FA for a specific user, make a DELETE request to this endpoint with the following JSON body: `{"user":"user name"}`.
   - The API disables 2FA for the specified user account.

Please refer to this documentation when working with 2FA in our system. Implementing 2FA significantly enhances the security of user accounts and provides an extra layer of protection against unauthorized access.

We have a http + javascript sample in user_test.go.
*/

// To test 2fa

// 1. Run mnms root `./mnmsctl/mnmsctl.exe -R -n root`
// 2. Add new user with API or web UI
// 3. Use following web page [generate secre] button to generate secret, and enable user's 2fa, scan QR code with google authenticator
//    to get 2fa code
// 4. Login user account which createed in step 2, mnms response sessionID instead of token, use this sessionID and 2fa code
//    To varify 2fa

// This is a sample code of frontend to demo how to generate qr code and validate code.
// To use this sample code, you need to add '*' in the mnms cors here
// 		r.Use(cors.Handler(cors.Options{
//		   AllowedOrigins:   []string{"https://*", "http://*", "*"},

/*
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>2FA Sample</title>
    <script src="https://cdn.jsdelivr.net/npm/qrcode@1.4.4/build/qrcode.min.js"></script>
    <style>
        .json-contianer {

            padding: 10px 0;
            font-family: monospace;
            white-space: pre;
        }

        div {
            margin: 10px 0;
        }
    </style>
</head>

<body>
    <!-- add login  -->
    <h2>Login</h2>
    <input type="text" id="username" placeholder="Enter username">
    <input type="password" id="password" placeholder="Enter password">
    <button id="login">Login</button>
    <div id="loginResult" class="json-contianer">...</div>

    <h2>Two-Factor Authentication (2FA) Sample</h2>
    <input type="text" id="username-2fa" placeholder="Enter username">
    <button id="generate">Generate Secret</button>
    <button id="update">Update Secret</button>
    <div id="2fa-operation-result"></div>
    <div id="">Session ID: <span id="sessionID"></span></div>

    <div id="qrcode"></div>

    <input type="text" id="code" placeholder="Enter 2FA Code">
    <button id="validate">Validate Code</button>
    <div id="validationResult"></div>

    <script>
        const generateBtn = document.getElementById('generate');
        const updateBtn = document.getElementById('update');
        const operationResult = document.getElementById('2fa-operation-result');


        const codeInput = document.getElementById('code');
        const loginButton = document.getElementById('login');
        const loginResult = document.getElementById('loginResult');
        const validateBtn = document.getElementById('validate');
        const validationResult = document.getElementById('validationResult');
        const qrCodeDiv = document.getElementById('qrcode');
        const sessionIDSpan = document.getElementById('sessionID');
        var token =""

        let secret = '';

        // login
        loginButton.addEventListener('click', async () => {
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const response = await fetch('http://localhost:27182/api/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user: username,
                    password: password,
                }),
            });

            const data = await response.text();
            const jsonObject = JSON.parse(data);
            if (jsonObject.token) {
                token = jsonObject.token;
            }
            // if sessionID is not empty, write to sessionIDSpan
            if (jsonObject.sessionID) {
                sessionIDSpan.textContent = jsonObject.sessionID;
            }
            const prettyJsonString = JSON.stringify(jsonObject, null, 2);// pretty print JSON

            loginResult.textContent = prettyJsonString;
        });
        // update 2fa
        updateBtn.addEventListener('click', async () => {
            account = document.getElementById('username-2fa').value;
            const response = await fetch('http://localhost:27182/api/v1/2fa/secret', {
                method: 'PUT',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user:account,
                }),
            });
            const data = await response.text();
            const jsonObject = JSON.parse(data);
            const prettyJsonString = JSON.stringify(jsonObject, null, 2);// pretty print JSON
            operationResult.textContent = prettyJsonString;
            secret = jsonObject.secret;
            showQRcode(jsonObject.account, jsonObject.secret);
        });


        // new 2fa
        generateBtn.addEventListener('click', async () => {
            account = document.getElementById('username-2fa').value;
            const response = await fetch('http://localhost:27182/api/v1/2fa/secret', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + token,
                },
                body: JSON.stringify({

                    user:account,
                }),
            });
            const data = await response.text();
             const jsonObject = JSON.parse(data);
            const prettyJsonString = JSON.stringify(jsonObject, null, 2);// pretty print JSON
            operationResult.textContent = prettyJsonString;
            secret = jsonObject.secret;
            showQRcode(jsonObject.account, jsonObject.secret);
        });

        function showQRcode(account,secret) {
            qrCodeDiv.innerHTML = '';
            const otpAuthURL = `otpauth://totp/Atop_MNMS:${account}?secret=${secret}&issuer=Atop_MNMS`;
            QRCode.toDataURL(otpAuthURL)
                .then((url) => {
                    let img = new Image();
                    img.src = url;
                    qrCodeDiv.appendChild(img);
                })
                .catch((err) => {
                    console.error(err);
                });
        }


        validateBtn.addEventListener('click', async () => {
            // get session id


            const code = codeInput.value;
            console.log("code : ",code, "session id: ", sessionIDSpan.textContent)
            const response = await fetch('http://localhost:27182/api/v1/2fa/validate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    sessionID: sessionIDSpan.textContent,
                    code: code,
                }),
            });

            const data = await response.text();
            validationResult.textContent = data;
        });

    </script>
</body>
</html>
*/
