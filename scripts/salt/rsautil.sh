#!/bin/sh
set -eu
OPENSSL="/usr/bin/openssl"
TMPDIR="/tmp/`whoami`"
mkdir -p ${TMPDIR}

TMP=${TMPDIR}/RSA_TMP_ENC
TMP1=${TMPDIR}/RSA_TMP_ENC_SIN
echo "$OPENSSL rsautl -encrypt -in $4 -inkey $2 -pubin -out $TMP"
echo "$OPENSSL rsautl -sign -in $TMP -inkey $3 -out $TMP1"
$OPENSSL rsautl -encrypt -in $4 -inkey $2 -pubin -out $TMP
$OPENSSL rsautl -sign -in $TMP -inkey $3 -out $TMP1
/usr/bin/hexdump -C $TMP1|awk 'BEGIN{out="";}{for(i=2;i<NF && i<=17;i++)out = out""$i}END{print out}' > $5