# MNMS authentication guide

## Custom mnms private key
`mnms`,  is designed with security in mind and uses encryption to protect data. The system comes with an RSA private key and its corresponding public key for encryption and decryption. However, users have the option to configure their own private key when running `mnms`.

```
# run mnms root with default private key
mnmsctl -R -n root
# run mnms root with custom private key
mnmsctl -R -n root -privkey {private_key_file}
```
We suggest use custom private key since it is more secure.

Be aware that for any reason, if the user restart `mnms` with different private key, the `config.json` file will not be correctly decrypted. Users who provide private key must keep it safe and use it consistently.


### Get mnms public key

`mnms` provides a utility to get the public key. The public key is generated from the private key used to start the `mnms`.

```bash
mnmsctl util -mnmspubkey -out {out_file}
```

You can specify the output file using the `-out` flag. If you run the `-mnmspubkey` flag without the `-out` flag, the output will be written to stdout. For example, you can obtain the `mnms` public key as follows:

```bash
./mnmsctl/mnmsctl util -mnmspubkey > ~/mnmspub
```

## Default administrator password
The default administrator password is "default". After logging in, you should change it.  
```
account: admin
password: default
```

### Login
User can login with API or Web UI. mnms will return a token after login. The token will be used for other API request.  

POST /api/v1/login  
```json
{
    "user": "admin",
    "password": "default"
}
```

## Change password
Should login with administrator permission first, then use the following API to change the password. Or you can use Web UI user management page to change the password also.

API need to set the `Authorization` header with the token returned by login API.

PUT /api/v1/users
```json
{
    "name": "admin",
    "password": "new_password",
    "email":"email@aa.bb.cc",
    "role": "admin"
}
```

## Password rule
Over 8 characters, at least one uppercase letter, one lowercase letter, one number and one special character. Allowed special characters are `@$!%*#?&`

### Generate RSA key pair

`mnms` can help users to create an RSA key pair. RSA keys are mainly used to encrypt data exported by `mnms`. `mnms` does not want to export unencrypted data directly, so when exporting, it will require the user to provide a public key for encryption. Users can use the corresponding private key to decrypt the content.

**Example**: default output

Running `nmnmsctl util -genrsa` without any flags will write a key pair to the default file locations: `$HOME/.mnms/id_rsa` and `$HOME/.mnms/id_rsa.pub`.

```bash
mnmsctl util -genrsa
```

**Example**: `-name` flag

To generate an RSA key pair with the `name` flag, use it as a prefix for the key pair file name. For example, if you use `name ~/mnmskey`, the private key will be stored in `~/mnmskey`, and the public key will be stored in `~/mnmskey.pub`.

```bash
mnmsctl util -genrsa -name ~/mnmskey
```




## Encrypt

The `mnms` tool also provides an encryption function. Users can choose any RSA tool to encrypt data. The following sample describes how to encrypt a file with `mnmsctl`.

```bash
mnmsctl util -encrypt -pubkey {pubkey_file} -in {input_file} -out {output_file}
```

If the `-in` or `-out` flag is not provided, stdin and stdout will be used instead.

## Decrypt

```bash
mnmsctl util -decrypt -privkey {private_key} -in {input_file} -out {output_file}
```

If the `-in` or `-out` flag is not provided, stdin and stdout will be used instead.



## Get `config.json`

`config.json` is the configuration file for MNMS. You can export an encrypted version of this file using the command `mnmsctl util -export -configfile -pubkey {pubkey_file} -privkey ~/.mnms/privkey`. The encryption is done using the public key that was set when running command, so you will need the corresponding private key to decrypt the file. if `-privkey` is not provided, the default private key will be used.

```bash
./mnmsctl/mnmsctl util -export -configfile -pubkey ~/pubkey.pem -privkey ~/privkey > ~/mnmsconfig
```

The above command will export the `config.json` file to the `~/mnmsconfig` file. The `config.json` file is encrypted using the public key in the `~/.mnms/pubkey.pem` file. The private key in the `~/.mnms/privkey` file is used to decrypt the original config.json file. Private key is same as the private key used to start the `mnms`. If `mnms` start without specify `-privkey` flag, the default private key will be used. In this case, you can omit the `-privkey` flag. `-pubkey` and `-privkey` do not need to be a pair. You can use any public key to encrypt the file, and use any private key to decrypt the output file.


### Update `config.json`

When a user has made changes to the **`config.json`** file and wishes to update the MNMS system, they must first encrypt the file using the public key. The default public key can be obtained using the command **`mnmsctl util -mnmspubkey -out {filename}`**. The following is an example of encrypting the modified **`config.json`**

The use command `mnmsctl util -import -configfile -in {filename}` to upload encrypted `config.json` file

Step:  
1. Encrypt the modified `config.json` file using the proper public key, the paired key corresponding to the private key used to start the `mnms` service. 
```
mnmsctl util -encrypt -pubkey {pubkey_file} -in {modified_config} -out {encrypted_config}
```
2. Upload the encrypted `config.json` file to MNMS
```
mnmsctl util -import -configfile -in {encrypted_config}
```


# Two-Factor Authentication (2FA) Guide

## What is 2FA?

Two-Factor Authentication (2FA) is an additional layer of security that requires users to provide two forms of identification before accessing an account. The first form of identification is typically a password, and the second form of identification is usually a unique code that is generated by an app or sent to a user's smart phone.

## Why use 2FA?

2FA adds an extra layer of protection to user accounts and can prevent unauthorized access, even if a password is compromised. By requiring two forms of identification, 2FA helps to verify that a user is who they claim to be.

## Enabling 2FA

To enable 2FA for your account, follow these general steps:

1. Log in to your account and navigate to the user management settings.
2. Enable 2FA for your account.
3. A QR code will be displayed. Use an authenticator app, such as Google Authenticator or Authy, to scan the QR code and set up 2FA.
4. Verify that your 2FA is working by logging out and logging back in.

## Disabling 2FA

If you need to disable 2FA for your account, follow these steps:

1. Log in to your account and navigate to the user management settings.
2. Disable 2FA for your account.
3. Verify that 2FA has been disabled by logging out and logging back in.

Note: Disabling 2FA removes an important security feature from your account. Make sure to only disable 2FA if it's absolutely necessary and take other steps to secure your account.


## Conclusion

2FA is an important security feature that can help protect your accounts from unauthorized access. By following these best practices and enabling 2FA for your account, you can help ensure the security of your accounts and personal information.






