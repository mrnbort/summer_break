# simple curl test of 2 endpoints assuming default port for Flask, adjust as needed
curl -X POST http://127.0.0.1:8080/transactions  -F "file=@testdata/data.csv"
echo
curl http://127.0.0.1:8080/report
