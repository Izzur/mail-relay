# Mail Relay

API to relay incoming HTTP call to send SMTP email

### Quick Start

```
go run main.go
```

### Build

```
go build
```

Example Body

```
{
    "personalizations": [
        {
            "to": [
                {
                    "email": "john.doe@example.com",
                    "name": "John Doe"
                }
            ],
            "subject": "Hello, World!"
        }
    ],
    "content": [
        {
            "type": "text/plain",
            "value": "Heya!"
        }
    ],
    "from": {
        "email": "sam.smith@example.com",
        "name": "Sam Smith"
    },
    "reply_to": {
        "email": "sam.smith@example.com",
        "name": "Sam Smith"
    }
}
```