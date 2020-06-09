# Training
Example server

To run this example server 

1. Clone this repo
2. cd /path/to/repo
3. docker-compose up
4. Server listens on port 9765

## API Documentation

**Note:** All balances are in ₹ * 100. E.g., balance of 300000 shows a balance of ₹3000.00. Maintaining balances in paise allows easy calculation of rounding.

### Account

**POST /accounts** - Create a new account

Inputs: 

userid string

balance int         *- This is the initial balance to create the account*

Returns: New account id

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

**GET /accounts/[accountid]** - Get the activity for an account

Inputs: None

Returns: List of transactions on the account

Returns an error if the account number doesn't exist 

Operations:

- OB: Original Balance
- CR: Credit
- DB: Debit

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
