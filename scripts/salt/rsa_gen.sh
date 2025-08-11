#!/bin/sh
set -eu
OPENSSL="/usr/bin/openssl"
TMPDIR="${1}"
mkdir -p ${TMPDIR}
cd $TMPDIR
RSAPEM=rsa_tmp
PRI=${2}.pem
PUB=${2}_pub.pem
$OPENSSL genrsa -out $RSAPEM  2048
$OPENSSL pkcs8 -topk8 -inform PEM -in $RSAPEM -outform PEM -nocrypt -out $PRI
$OPENSSL rsa -in $PRI -pubout -out $PUB
