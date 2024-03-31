package main

import (
	"github.com/samborkent/uuidv7"
	"time"
	"fmt"
	"encoding/binary"
)

func main() {


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
	fmt.Printf("%s|%s|%d|%s|%s\n",str,shortStr,sequenceNumber,timestrUtc,timestrLocal)

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
