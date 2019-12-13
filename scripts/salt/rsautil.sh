#!/bin/sh
set -eu
OPENSSL="/usr/bin/openssl"
TMPDIR="/tmp/`whoami`"
mkdir -p ${TMPDIR}

if [ "$#" -lt 2 ]
then
	echo "usage 1 ./rsautil enc app_pub_key sys_pri_key inputfile outputfile"
	echo "usage 2 ./rsautil dec sys_pub_key app_pri_key inputfile outputfile"
	echo "usage 3 ./rsautil genappkey app"
	echo "usage 4 ./rsautil gensyskey sys"
	exit
fi

if [ "$1" = "enc" ]
then
	TMP=${TMPDIR}/RSA_TMP_ENC
	TMP1=${TMPDIR}/RSA_TMP_ENC_SIN
    echo "$OPENSSL rsautl -encrypt -in $4 -inkey $2 -pubin -out $TMP"
    echo "$OPENSSL rsautl -sign -in $TMP -inkey $3 -out $TMP1"
	$OPENSSL rsautl -encrypt -in $4 -inkey $2 -pubin -out $TMP
	$OPENSSL rsautl -sign -in $TMP -inkey $3 -out $TMP1
	/usr/bin/hexdump -C $TMP1|awk 'BEGIN{out="";}{for(i=2;i<NF && i<=17;i++)out = out""$i}END{print out}' > $5
elif [ "$1" = "dec" ]
then
	TMP=${TMPDIR}/RSA_TMP_BIN
	TMP1=${TMPDIR}/RSA_TMP_CHK
	./hex2bin $4 $TMP
    echo "$OPENSSL rsautl -verify -in $TMP -inkey $2 -pubin -out $TMP1"
    echo "$OPENSSL rsautl -decrypt -in $TMP1 -inkey $3 -out $5"
	$OPENSSL rsautl -verify -in $TMP -inkey $2 -pubin -out $TMP1
	$OPENSSL rsautl -decrypt -in $TMP1 -inkey $3 -out $5
elif [ "$1" = "genappkey" ]
then
    TMP=${TMPDIR}/rsa_tmp
    PRI=${2}.pem
    PUB=${2}_pub.pem
	$OPENSSL genrsa -out $TMP  2048
    $OPENSSL pkcs8 -topk8 -inform PEM -in $TMP -outform PEM -nocrypt -out $PRI
    $OPENSSL rsa -in $PRI -pubout -out $PUB
elif [ "$1" = "gensyskey" ]
then
    TMP=${TMPDIR}/rsa_tmp
    PRI=${2}.pem
    PUB=${2}_pub.pem
	$OPENSSL genrsa -out $TMP  4096
    $OPENSSL pkcs8 -topk8 -inform PEM -in $TMP -outform PEM -nocrypt -out $PRI
    $OPENSSL rsa -in $PRI -pubout -out $PUB
fi
