# Go uuidv7 version

[Go uuidv7 module](https://github.com/samborkent/uuid)

```bash
go mod init uuidv7
go mod tidy
go build

# example run
uuidv7
018e9508-4461-7149-b405-03f6a362ce41|03f6a362ce41|329|2024-03-31 15:02:10.785 +0000 UTC|2024-03-31 18:02:10.785 +0300 EEST
```

Return five values, delimiter |
* uuidv7
* short 
* sequence
* timestamp UTC
* timestamp local

[More documentation](github.com/samborkent/uuidv7)

## Using uuidv7 with Bash / Ksh /.... 
Look uuidv7.sh
```bash
IFS="|" read uuidv7 short seq timestampUTC timestampLocal xstr < <(uuidv7)
echo "uuidv7:$uuidv7"
```

