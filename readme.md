# Go learning

The purpose of this repository is to track my progress in learning Go. I will create an simple IDP (Identity Provider) using Go. The IDP need to have the following features:

- [x] Create application
- [x] User registration
- [ ] User login
- [ ] User logout
- [ ] Authorization

## Create application

An application is a client that will use the IDP to authenticate users. The application will have a name and a secret key. The secret key will be used to authenticate the application with the IDP. This secret key will be generated by the IDP and sent to the application when it is created. The IDP will generate a passphrase that will be used to reset the secret key in case it is lost or compromised.

### Parameters

- JWT_EXPIRATION: The expiration time of the JWT token.
- JWT_EXPIRATION_REFRESH: The expiration time of the JWT refresh token.

## User registration

A user can register by providing a username and password. The request should come with an API key.