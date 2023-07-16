# Summer Break Service

## Description

The Summer Break Service collects transaction data of expenses and 
income and provides a report with gross revenue, expenses, and net 
revenue.

## Run in Docker

1. Copy docker-compose.yml

   - change the ports if needed

2. Start a container with `docker-compose up`

## Architectural Summary

The Summer Break Service has two public API endpoints: POST 
/transactions and GET /report. To make requests, the service uses 
injected `http.Server`. Read Timeout is set to 5 seconds by default 
but can be changed from the command line. Write Timeout is set to 
30 seconds by default but can be changed from the command line.

POST handler reads every line in the file with `csv.NewReader` 
and `.Read()` while skipping lines with the wrong number of fields. 
After the lines from the file are parsed, they are passed to the 
`ProcessTransactions()` function which saves the data in the server's 
local cache. This approach allows on-the-fly processing without the 
need to load everything to the memory which can be beneficial in 
case of large input files.

GET handler generates the revenue/expenses report by calling the 
`GenerateReport()` function. This function calculates the gross revenue
by summing all "Income" transactions, expenses by summing all 
"Expense" transactions, and finally net revenue by subtracting 
expenses from the gross revenue. It returns a JSON document with the 
three values.

## Error handling

The service is designed to be as resilient as possible. It won't 
fail on non-critical issues, like one of the endpoints 
non-responsive, or data parsing failed. These errors are logged 
to stdout and the service continues to work. The only critical 
errors are when the service is unable to parse the command line parameter or environment
variables, or if the server fails to run. In such cases the service 
exits with non-zero exit code.

The service also handling both `SIGINT` and `SIGTERM` and 
exits gracefully. It is done by canceling context. This context 
is passed to all the functions.

## Dependency injection

In this project, dependency injection is achieved through 
"constructor" functions and by directly setting fields of the 
struct. All the necessary configuration/wiring is done in `main.go`. 

## Testing

The code of this project covered with unit and integration tests. 
For mocking the indispensable `moq` code generator is used.
All the tests passed with `-race` flag as well. 

## Code style and comments

The code is formatted with `gofmt` and `goimports`. Comments are written
in the GoDoc style and are used to explain the logic that may not be 
obvious from the code.

## API

API endpoints:

1. `POST /transactions` - takes as input CSV formatted data as below,
parses it, and saves it.
```
2020-07-01, Expense, 18.77, Gas
2020-07-04, Income, 40.00, 347 Woodrow
2020-07-06, Income, 35.00, 219 Pleasant
2020-07-12, Expense, 49.50, Repairs
```
- Returns:
```json
{
  "status": "ok" 
}
```
- Example of usage:
```
curl -X POST http://127.0.0.1:8080/transactions  -F "file=@testdata/data.csv"
```

2. `GET /report` - return a JSON document with the tally of gross 
revenue, expenses, and net revenue (gross - expenses) as follows:
```json
{
    "gross-revenue": 0,
    "expenses": 0,
    "net-revenue": 0
}
```
- Example of usage:
```
curl http://127.0.0.1:8080/report
```

## General considerations

I made an assumption for this service that it is acceptable to lose the 
collected data if server's process failed. This is not what we would like
to have in the real project; however for such a toy example it should be fine.

One of the shortcomings of the proposed solution is that if the server
shuts down, all the transaction data will be lost and when the server 
is restarted, it will need all the previously processed CSV files to be 
sent again. This can be solved by connecting the service to a database 
where the transaction data will be saved permanently. 

Another potential issue is that the endpoints are not protected, so 
anyone who knows the URLs can access sensitive financial information of 
the main user.

Model uses float64 for money-related concern which is fine for 
a toy example, but in a real project, we would need a more appropriate
decimal type.

These limitations are intentional due to strict time restrictions.

## How to run the binary directly

```
Usage:
  summer_break [OPTIONS]

Application Options:
      --port=                  http data server port (default: 8080)
      --http-read-timeout=     timeout for read HTTP requests (default: 5s)
      --http-write-timeout=    timeout for write HTTP requests (default: 30s)

Help Options:
  -h, --help            Show this help message
```

## Potential improvements

1. Introduce persistent storage to keep all the transactions in.
2. Add an authentication method to limit access
3. Add a rate limiter to prevent attacks.
4. Use a more appropriate decimal type for money handling operations.
5. Add validation for transaction type to validate that it is either 
"Expense" or "Income".
6. Add validation for reasonable amount values (positive only, 
no greater than X amount).
7. Limit maximum body size for a POST request to prevent attacks.