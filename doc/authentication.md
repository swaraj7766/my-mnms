# MNMS authentication guide

# First Run

`mnms` creates a default admin account with the username "admin". Users can retrieve the password for this account using the `/api/vi/getpass` API.

```bash
curl -X POST -H 'Accept: application/json'  https://mnms.atop.com/api/v1/getpass -d '{"email":"admin","password":"temporary_password"}'
```

The response contains a warning because `admin` is not a valid email address; simply ignore it.

```json
{
    "status": "warning: 553 5.1.3 The recipient address <admin> is not a valid RFC-5321 address. Learn\n5.1.3 more at\n5.1.3  https://support.google.com/mail/answer/6596 d22-20020a9d5e16000000b006864b5f4650sm6157422oti.46 - gsmtp",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluIiwiZXhwIjoxNjc1NzUzNjUxLCJwYXNzIjoiMTIzNCJ9.OMbKg2RBIQcmhT9c-hA3Xho_9Dx7F8OHIuZox2AW7Sw"
}
```

POST to `/api/v1/pw/{token}` with body `{"email":"admin", "password":"temporary_password"}` to retrieve the admin's real password. The {token} field is from the  `/api/vi/getpass`  request’s response. Note that the URL is only valid for 10 minutes.  `mnms` should response real password of admin.

Except for the administrator, the first email to make the `/api/v1/getpass` request will be treated as the admin. This email does not need to exist in the system. Any other email must be added to `mnms` first by sending a POST request to `/api/v1/users` with a body of `{"email","is_admin"}`

## How to retrieve password

1. Add an email to `mnms` by making a POST request to `/api/v1/users` with a body of `{"email": <EMAIL>, "is_admin": <BOOLEAN>}`. `is_admin` indicates whether the added email has admin permissions.
2. POST to `/api/v1/getpass` with a body of `{"email":<EMAIL>, "password":<PASSWORD>}` to generate a temporary token for retrieving the real password. `mnms` will send a link to the email provided previously. Go to the link, answer the temporary password, and retrieve the real password.