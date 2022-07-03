#!/bin/bash
NAME=${1}
ACTUAL=${2}
ANSWER=${3}

if [ ${ACTUAL} = ${ANSWER} ]; then
echo "[Pass] ${NAME}"
else
echo "[Failed] ${NAME} : output = ${ACTUAL} , answer = ${ANSWER}"
exit 1
fi
