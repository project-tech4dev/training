# Training

To run this example server 

1. Clone this repo
2. cd /path/to/repo
3. docker-compose up
4. Server listens on port 9765

If you get a new pull from the Repo, run the following commands to rebuild and restart the docker services

```json
> docker-compose down
> docker-compose up —build
```

## API Documentation

**Note:** All balances are in ₹ * 100. E.g., balance of 300000 shows a balance of ₹3000.00. Maintaining balances in paise allows easy calculation of rounding.

Authorization header is required for operations that require authorization. The format of the Auth header is 

```json
Authorization: Bearer <Token>
```

Bank Manager already has an account in the system (id: bankmanager, pwd: headhoncho). There are two roles in the system BankManager and User. The Authorization notes below say what roles are permitted to do that operation.

### Account

**POST /accounts** - Create a new account

Inputs: 

userid string

balance int         *- This is the initial balance to create the account*

Returns: New account id

Authorization: BankManager

Example 

```json
POST /accounts
{
	"userid" : "tester",
	"balance": 600000
}

Response
{
  "accountid": "10005",
	"balance": 600000
}
```

**POST /accounts/credit** - Deposit money into the account

Inputs: (user id and initial balance)

accountid string

amount int

Returns: New balance

Returns an error if the account number doesn't exist 

Authorization: User (account owner)

Example 

```json
POST /accounts
{
	"accountid": "10005",
	"amount": 200000
}
Response
{
  "balance": 800000
}
```

**POST /accounts/debit** - Withdraws money from the account

Inputs: 

accountid string

amount int

Returns: New balance

Returns an error if the account number doesn't exist or if there aren't sufficient funds

Authorization: User (account owner)

Example 

```json
POST /accounts/debit
{
	"accountid": "10005",
	"amount": 100000
}
Response
{
  "balance": 700000
}
```

**GET /accounts/`<accountid>`** - Get the activity for an account

Inputs: None

Returns: List of transactions on the account

Returns an error if the account number doesn't exist 

Operations:

- OB: Original Balance
- CR: Credit
- DB: Debit

Authorization: User (account owner) or BankManager

Example 

```json
GET /accounts/10005

Response
{
  "accountid": "10005",
  "activity": [
    {
      "amount": 600000,
      "operation": "OB",
      "balance": 600000,
      "createdat": "2020-06-09T00:52:47.541942Z"
    },
    {
      "amount": 200000,
      "operation": "CR",
      "balance": 800000,
      "createdat": "2020-06-09T01:03:09.5493978Z"
    },
    {
      "amount": 100000,
      "operation": "DB",
      "balance": 700000,
      "createdat": "2020-06-09T01:05:13.2286343Z"
    }
  ]
}

```

**GET /accounts** - Deposit money into the account

Inputs: None

Returns: All the accounts with all their activity

Authorization: BankManager

Example 

```json
GET /accounts

Response
{
  "accounts": [
    {
      "id": "10002",
      "userid": "tester",
      "createdat": "2020-06-09T00:01:32.504121Z",
      "balance": 450000,
      "activity": [
        {
          "amount": 450000,
          "operation": "OB",
          "balance": 450000,
          "createdat": "2020-06-09T00:01:32.504113Z"
        }
      ]
    },
    {
      "id": "10003",
      "userid": "tester",
      "createdat": "2020-06-09T00:02:04.819345Z",
      "balance": 700000,
      "activity": [
        {
          "amount": 500000,
          "operation": "OB",
          "balance": 500000,
          "createdat": "2020-06-09T00:02:04.819327Z"
        },
        {
          "amount": 200000,
          "operation": "CR",
          "balance": 700000,
          "createdat": "2020-06-09T00:05:41.5509829Z"
        }
      ]
    },
    {
      "id": "10004",
      "userid": "tester",
      "createdat": "2020-06-09T00:05:57.7149467Z",
      "balance": 600000,
      "activity": [
        {
          "amount": 600000,
          "operation": "OB",
          "balance": 600000,
          "createdat": "2020-06-09T00:05:57.7149255Z"
        }
      ]
    },
    {
      "id": "10005",
      "userid": "tester",
      "createdat": "2020-06-09T00:52:47.5419653Z",
      "balance": 700000,
      "activity": [
        {
          "amount": 600000,
          "operation": "OB",
          "balance": 600000,
          "createdat": "2020-06-09T00:52:47.541942Z"
        },
        {
          "amount": 200000,
          "operation": "CR",
          "balance": 800000,
          "createdat": "2020-06-09T01:03:09.5493978Z"
        },
        {
          "amount": 100000,
          "operation": "DB",
          "balance": 700000,
          "createdat": "2020-06-09T01:05:13.2286343Z"
        }
      ]
    }
  ]
}
```

### Users

**POST /users -** Create User

Inputs: 

fullname string

username string

password string

Returns: ID of the new user

Returns an error if the user name is taken

Authorization: None

Example 

```json
POST /users
{
	"fullname": "O. B. B. G",
	"username": "obbg",
	"password": "ppp"
}

Response
{
  "userid": "1000004"
}

```

**GET /users** - Get a list of all users

Inputs: None

Returns: List of users (id, name, username, createdat, number of accounts) 

Authorization: BankManager

Example 

```json
GET /users

Response
{
  "users": [
    {
      "id": "1000002",
      "fullname": "My Name",
      "username": "a1003",
      "createdat": "11 Jun 20 19:10 +0000",
      "numaccounts": 5
    },
    {
      "id": "1000003",
      "fullname": "Second User",
      "username": "a1004",
      "createdat": "11 Jun 20 19:30 +0000",
      "numaccounts": 2
    },
    {
      "id": "1000004",
      "fullname": "O. B. B. G",
      "username": "obbg",
      "createdat": "11 Jun 20 20:30 +0000",
      "numaccounts": 0
    }
  ]
}
```

**GET /user/`<ID>`** - Get User

Inputs: None

Returns: The user including accounts and account activity 

Returns an error if the user name or password doesn't match

Authorization: User (self) or BankManager

Example 

```json
GET /user/1000004

Response
{
  "user": {
    "userid": "1000004",
    "name": "O. B. B. G",
    "username": "obbg",
    "password": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "role": "User",
    "createdat": "2020-06-11T20:30:18.564486Z",
    "accounts": []
  }
}
```

**POST /login** - Login

Inputs: 

username string

password string

Returns: ID of the new user and Auth token

Returns an error if the user name or password doesn't match

Authorization: None

Example 

```json
POST /login
{
	"username": "obbg",
	"password": "ppp"
}

Response
{
  "id": "1000004",
  "token": "c4289629b08bc4d61411aaa6d6d4a0c3c5f8c1e848e282976e29b6bed5aeedc7"
}
```

