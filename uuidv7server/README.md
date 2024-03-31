# Microservice, return uuidv7
## Use Go-language

 * [Go uuidv7 module](https://github.com/samborkent/uuid)
 * [uuidv7 command](https://github.com/kshji/go/tree/master/uuid7)


```bash
go mod init uuidv7server
go mod tidy
go build

# example run
uuidv7server -p 8888 -l 127.0.0.1
```

Test it:
```bash
wget -q -O - "http://localhost:8888/uuidv7"

curl -X GET "http://localhost:8888/uuidv7"
```

Return five values, delimiter |
 * uuidv7
 * short
 * sequence
 * timestamp UTC
 * timestamp local

[More documentation](github.com/samborkent/uuidv7)

