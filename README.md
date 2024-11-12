## Quick Start

To start:

```
docker-compose up -d
```

To stop and remove everything:

```
docker-compose down -v
```

On startup database will be seeded with 5 users with ids from 1 to 5 and 0 balance.

## Endpoints

### Get balance

```
GET http://localhost:8080/user/{userId}/balance
```

*Response:*

```
{
    "userId": 3,
    "balance": "10.00"
}
```

*Error response:*

```
{
    "errorText": "failed to get balance"
}
```


### Update balance

```
POST http://localhost:8080/user/{userId}/transaction
```

*Required header:* `Source-Type`

Possible values: `game`, `server`, `payment`

*Request body:*

```
{
    "state": "win",
    "amount": "10.00",
    "transactionId": "111-222-333"
}
```

* state - required, possible values `win` or `lose`

* amount - required, positive number

* transactionId - required, anything unique

Endpoint returns 200 code if transaction was successful and error response otherwise.

*Error response:*

```
{
    "errorText": "failed to update balance"
}
```
