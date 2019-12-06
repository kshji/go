package main

// parse_env.go
// example how to make template including "variables" and expand those
// - run this example:
// go run parse_env.go
//
// more about maps/hash
// https://yourbasic.org/golang/maps-explained/

import (
        "fmt"
        "regexp"
	"os"
	"strings"
)

// variable map - hash table

func main() {
	// parse all env variables using indexed array - map
        var varstr string = ""
	var envvars = make(map[string]string)

	// loop all environment variables and add to the map - hash table
        for _, e := range os.Environ() {
                pair := strings.SplitN(e, "=", 2)
                fmt.Println(e," - ",pair[0],"=",pair[1])
                varstr = "{"+pair[0]+"}"  // add {  } = string comparing is easier, because template using {}
                envvars[varstr] = pair[1]
        }

	// test template
        input := `
Expand env variables like 
- HOME:{HOME} 
or 
user&term&user:{USER}_{TERM}_{USER}  multivariable and some of those are multi

Not used variable ${FOOFOO} is not expanded

Path:
{PATH}
etc.
`
	// regexp rule to find {XXX} strings
	var r *regexp.Regexp
	r = regexp.MustCompile(`{([^}]+)}`)
	result := Expand(input,r,envvars)
	fmt.Println("-------------------------------")
	fmt.Println("Input:",input)
	fmt.Println("-------------------------------")
	fmt.Println("Result:",result)
	fmt.Println("-------------------------------")
	input="test{HOME}"
	fmt.Println(input,":",Expand(input,r,envvars) )
	input="test{HOME}{HOME}"
	fmt.Println(input,":",Expand(input,r,envvars) )
	input="test{HOME}{USER}"
	fmt.Println(input,":",Expand(input,r,envvars) )

}

func Expand(instr string, regexprule *regexp.Regexp, variables  map[string]string )  string {
	// input include 0-n {variables} and variables are of course not allways the same
	// loop all matches from input match by match and make unique replace all of those
        result:=regexprule.ReplaceAllStringFunc(instr,
		func(inputstr string) string {
			// search if there is {variable} ..., find 1st match
			match := regexprule.FindString(inputstr)
			if  match == "" { return inputstr }	// works without this, but faster+nicer 
			// replace simple use map/has table, don't use pointer, because map is already pointer to the hmap
			// little confusing but on this way we are allways sure that you use pointers :)
			// in this case you can make this easy, but here is only in example purpose used function
			//replace := envvars[match]
			replace := expandEnv(variables,match)
			if  replace == "" { return inputstr }
			// replace it
			return regexprule.ReplaceAllString(inputstr, replace)
		})
	return result
}


// expand variable name and look if variable exists, return value of variable
func expandEnv(keys map[string]string, key string) string {
	// look key from keys mapped/indexed hash table/array
        value:=keys[key]    // - search value using key
        return value
}

