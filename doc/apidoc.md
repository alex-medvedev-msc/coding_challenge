# API documentation

## Get Accounts

``GET /v1/accounts``

Returns all accounts from db without any filtering or pagination

**Responce example:**

```json
[
    {
        "id": "string_id",
        "owner": "some_user",
        "balance": 1.323,
        "currency": "PHP"
    }
]
```

**Possible errors:**

 - Status 500: internal server error


## Get Payments

``GET /v1/payments``

Returns all payments from db without any filtering or pagination

**Responce example:**

```json
[
    {
        "account": "bob123",
        "amount": 100,
        "to_account": "alice456",
        "direction": "outgoing"
    },
    {
        "account": "alice456",
        "amount": 100,
        "from_account": "bob123",
        "direction": "incoming"
    }
]
```

**Possible errors:**

 - Status 500: internal server error

## Transfer funds

``POST /v1/payments``

Transfer funds update balances of two corresponding accounts
and creates two matching payments with opposite direction in db:
incoming and outgoing with same amount.

Accounts must have same currency to do funds transfer
Sender must have enough funds to transfer them

**Request example:**

```json
{
    "from_account": "bob123",
    "amount": 100,
    "to_account": "alice456"
}
```

**Response example:**

Just 200 OK

**Possible errors:**

 - Status 500: internal server error
 - Status 404: one of accounts was not found
 - Status 409: there is currency mismatch or not enough funds
 - Status 401: malformed request, one of the account ids is empty or amount is negative

