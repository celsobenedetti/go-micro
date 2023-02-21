POST http://localhost:8081/authenticate
{
    "email": "admin@example.com",
    "password": "verysecret"
}

HTTP 202
[Asserts]
body contains "error\":false"
body contains "admin@example.com"


POST http://localhost:8081/authenticate
{
    "email": "admin@example.com",
    "password": "wrongpassword"
}

HTTP 401
[Asserts]
body contains "Invalid credentials"
