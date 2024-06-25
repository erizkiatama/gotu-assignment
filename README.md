# GoTu Take Home Assignment

## How to Run
1. Make sure you have `docker` & `docker-compose` installed / running on your computer.
2. Run `docker-compose up --build` on the terminal
3. The API is ready to use



# API Docs

## User

### Create User (Register)

- URL: **localhost:8080/api/v1/user/register**
- Method: **POST**

#### Request
```
{
    "email": "test@example.com",
    "password": "password",
    "name": "Example Test"
}
```

#### Response
```
{
    "result": {
        "access": "access_token",
        "refresh": "refresh_token"
    }
}
```

### Login

- URL: **localhost:8080/api/v1/user/login**
- Method: **POST**


#### Request
```
{
    "email": "test@example.com",
    "password": "password"
}
```

#### Response
```
{
    "result": {
        "access": "access_token",
        "refresh": "refresh_token"
    }
}
```

## Book

### List of Books

- URL: **localhost:8080/api/v1/book**
- Method: **GET**

#### Request

#### Response
```
{
    "result": [
        {
            "id": 1,
            "title": "The Great Gatsby",
            "author": "F. Scott Fitzgerald",
            "description": "Description for The Great Gatsby",
            "price": 100000
        },
        {
            "id": 2,
            "title": "It Ends with Us",
            "author": "Colleen Hoover",
            "description": "Description for It Ends with Us",
            "price": 150000
        },
        ....
    ]
}
```

## Order

### Create Order

- URL: **localhost:8080/api/v1/order**
- Method: **POST**

#### Header
```
{
    "Authorization" : "Bearer {{access_token}}"
}
```

#### Request
```
{
    "details": [
        {
            "book_id": 1,
            "quantity": 1,
            "price": 100000
        },
        {
            "book_id": 2,
            "quantity": 3,
            "price": 450000
        }
    ]
}
```

#### Response
```
{
    "result": {
        "id": 1,
        "user_id": 1,
        "total_quantity": 4,
        "total_price": 550000,
        "details": [
            {
                "id": 1,
                "book_id": 1,
                "book_title": "",
                "quantity": 1,
                "price": 100000
            },
            {
                "id": 2,
                "book_id": 2,
                "book_title": "",
                "quantity": 3,
                "price": 450000
            }
        ]
    }
}
```

### List Order

- URL: **localhost:8080/api/v1/order**
- Method: **GET**

#### Header
```
{
    "Authorization" : "Bearer {{access_token}}"
}
```

#### Request

#### Response
```
{
    "result": [
        {
            "id": 1,
            "user_id": 1,
            "total_quantity": 4,
            "total_price": 550000,
            "details":[]
        },
        {
            "id": 2,
            "user_id": 1,
            "total_quantity": 1,
            "total_price": 100000,
            "details":[]
        },
    ]
}
```

### Detail Order

- URL: **localhost:8080/api/v1/order/:orderId**
- Method: **GET**

#### Header
```
{
    "Authorization" : "Bearer {{access_token}}"
}
```

#### Request

#### Response
```
{
    "result": {
        "id": 1,
        "user_id": 1,
        "total_quantity": 4,
        "total_price": 550000,
        "details": [
            {
                "id": 1,
                "book_id": 1,
                "quantity": 1,
                "price": 100000
            },
            {
                "id": 2,
                "book_id": 2,
                "quantity": 3,
                "price": 450000
            }
        ]
    }
}
```
