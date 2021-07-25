# User SignUp and Reset Password

Implement simple user sign up and reset password functionality in Golang and MongoDB


## System Requirements

```bash
Go (Version 1.16)
MongoDB (Version 4.4 above)
```

## Installtion

Below command will Install all the dependencies recursively.

```bash
go get -d ./...
```

## Starting the MongoDB service

```bash
brew services start mongodb-community
```

## Starting the GO server

Use the below command to create executable and the run executable.

```bash
go build
./golang_mongodb
```

## Further Improvement's TODO

1. Seperate helper module or folder should be created for common functions which needs to be used all across repo.
2. Seperate config file or module should be created to manage all configs.
3. Folder structure should be improved to manange the code complexity.
