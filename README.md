# funds-loader
    
An application that loads money on customer accounts.

## Assumptions

I work with customer money as with float64. For production development, I would use something more appropriate package like this one https://github.com/Rhymond/go-money which doesn't have a problem with rounding.

I have +1 line comparing to the output.txt because of a transaction with ID `"id":"6928"` which appears twice in input.txt, but only once in output.txt


## How to build and run an application?

To build an application you need to execute `make build`
After that executed file will be stored in `bin` folder


To run an application you need to execute `./bin/funds-loader` or `go run main.go`
Results can be found `./testdata/result.txt`

## How to run unit tests?

To run unit tests you need to execute `make test`
