# CRUD

- https://github.com/mongodb/mongo-go-driver
- https://docs.mongodb.com/drivers/go/current/

```sh
# create a person
$ curl -X POST http://localhost:8080/person -H 'Content-type: application/json' -d '{"firstname":"hello","lastname":"world!"}'

# get all people
$ curl -X GET localhost:8080/person

# get one person
$ curl -X GET localhost:8080/person/61b04302fc6054c46bb2e8d9
```
