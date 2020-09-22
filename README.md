## About:
Simple web app implementing Rest API to save and send email messages. Messages that are older than 5 minutes are deleted.

## Requirements:
* Go v1.14.2 (or higher)
* Cassandra v3.11.7 (or higher)
* Mailgun account

## Setup:
Before running, you need to setup your Cassandra keyspace. To do this, launch 'cqlsh' and type:

```sh
CREATE KEYSPACE message_sender_rest_api with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
CREATE TABLE message_sender_rest_api.messages(id UUID, email text, title text, content text, magic_number int, created_at timestamp, PRIMARY KEY(id));
CREATE INDEX ON message_sender_rest_api.messages(magic_number);
CREATE INDEX ON message_sender_rest_api.messages(email);
```

You also need to provide '.config' file with your Mailgun API key and domain name, like this:
```json
{
    "mailgun": {
        "api_key": "TYPE YOUR API KEY HERE",
        "domain": "type-your-domain-here.com"
    }
}
```

## Running:
```sh
go run *.go
```

## Usage examples:

Save a message:
```sh
curl -X POST localhost:8080/api/message -d '{"email":"john.doe@example.com","title":"hello there","content":"how are you john?","magic_number:123"}'
```

Send all messages with given 'magic_number':
```sh
curl -X POST localhost:8080/api/send -d '{"magic_number":123}'
```

Get all messages addressed to given email address:
```sh
curl -X GET localhost:8080/api/messages/{email}
```
