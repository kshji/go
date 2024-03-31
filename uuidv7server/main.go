package main

import (
	"github.com/samborkent/uuidv7"
	"time"
	"fmt"
	"encoding/binary"
	"net/http"
	"flag"
        "log"
        _"encoding/json"
        _"bytes"
        _"html"
        _"strings"
)

func main() {

	var port string
	var hostaddress string

	flag.StringVar(&port, "p", "3000", "port to listen on")
	flag.StringVar(&hostaddress, "l", "0.0.0.0", "ip to listen on")

	flag.Parse()
	address := hostaddress + ":" + port

	mux := http.NewServeMux()
	mux.HandleFunc("GET /uuidv7/", generateUuidv7)
	mux.HandleFunc("GET /", Hello)

	logger :=  logRequest(mux)

        // Start our HTTP server
        log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
        log.Println("Http uuidv7 server listen:", address)
        if err := http.ListenAndServe(address, logger); err != nil {
                log.Fatalln("Can't start, error:", err)
        }

}

func Hello(w http.ResponseWriter, r *http.Request) {
        log.Println("Home:")
        w.Write([]byte("Hello World!\nRead the documentation\n"))
}

func generateUuidv7(w http.ResponseWriter, r *http.Request) {
        log.Println("gen uuidv7:")

	var creationTimeBits [8]byte
	uniqueID := uuidv7.New()
	str := uniqueID.String()
	shortStr := uniqueID.Short()
	sequenceNumber := uniqueID.SequenceNumber()


	copy(creationTimeBits[:], uniqueID[:8])

	// Right shift timestamp bytes
	rightShiftTimestamp(creationTimeBits[:])

	//creationTime := uniqueID.CreationTime()
	timestrUtc := time.UnixMilli(int64(binary.BigEndian.Uint64(creationTimeBits[:]))).UTC()
	timestrLocal := time.UnixMilli(int64(binary.BigEndian.Uint64(creationTimeBits[:])))
	fmt.Fprintf(w,"%s|%s|%d|%s|%s\n",str,shortStr,sequenceNumber,timestrUtc,timestrLocal)

}

// Right shift timestamp bytes
func rightShiftTimestamp(uuid []byte) {
	uuid[7] = uuid[5]
	uuid[6] = uuid[4]
	uuid[5] = uuid[3]
	uuid[4] = uuid[2]
	uuid[3] = uuid[1]
	uuid[2] = uuid[0]
	uuid[1] = 0
	uuid[0] = 0
}


func logRequest(handler http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                log.Printf("%s|%s|%s|\n", r.RemoteAddr, r.Method, r.URL )
                for k, v := range r.Header {
                        log.Printf("  Header field %q, Value %q\n", k, v)
                }
                log.Println("ServeHTTP:")
                handler.ServeHTTP(w, r)
        })
}

