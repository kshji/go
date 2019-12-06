# Go Language - Parse OS environment variables and expand values in template #

This **parse_env.go** read all your environment variables to the map (hash table).

And also expand variable values in string template.

Example
```go
input := `
Hello world, my home dir is {HOME} and logname is {LOGNAME}. 
My terminal type is {TERM}
`
``` 

Result is:
```go
Hello world, my home dir is /home/myusername and logname is myusername.
My terminal type is xterm-256
```



