#!/bin/bash
echo "---------------"

echo "1. 200 check"
ACTUAL=`curl "http://localhost:8080/" -o /dev/null -w '%{http_code}\n' -s `
test/verify.sh "1-1. http://localhost:8080/ " "${ACTUAL}" "200" || exit 1

ACTUAL=`curl -u test:test "http://localhost:8080/summary/year/2021" -o /dev/null -w '%{http_code}\n' -s `
test/verify.sh "1-2. http://localhost:8080/summary/year/2021" "${ACTUAL}" "200" || exit 1

ACTUAL=`curl -u test:test "http://localhost:8080/record/recent/" -o /dev/null -w '%{http_code}\n' -s `
test/verify.sh "1-3. http://localhost:8080/record/recent/" "${ACTUAL}" "200" || exit 1

echo "2. Add record"
ACTUAL=`curl -u test:test -X POST -H "Content-Type: application/json" -d '{"category_id": 210, "price":1}' http://localhost:8080/record/ -o /dev/null -w '%{http_code}\n' -s `
test/verify.sh "2-1. Add nowtime record" "${ACTUAL}" "201" || exit 1

ACTUAL=`curl -u test:test -X POST -H "Content-Type: application/json" -d '{"category_id": 220, "price":1234, "memo":"test","date":"19990101"}' http://localhost:8080/record/ -o /dev/null -w '%{http_code}\n' -s `
test/verify.sh "2-2. Add oldtime record" "${ACTUAL}" "201" || exit 1

echo "3. Delete record"
ACTUAL=`curl -u test:test -X DELETE http://localhost:8080/record/1 -o /dev/null -w '%{http_code}\n' -s `
test/verify.sh "3-1. Delete 2-1 record" "${ACTUAL}" "204" || exit 1

echo "All test passed"
echo "---------------"
