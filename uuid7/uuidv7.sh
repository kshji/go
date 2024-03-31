#!/usr/bin/ksh
# or bash

IFS="|" read uuidv7 short seq timestampUTC timestampLocal xstr < <(uuidv7)
echo "uuidv7:$uuidv7"
