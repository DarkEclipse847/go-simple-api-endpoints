### go-simple-api

* To handle concurrency and to prevent "50X" Err, i used transactions and 'REPEATABLE READ' isolation-type
* All api endpoints were tested using postman
* Tests for all nesessary funcs

Run:  

``` bash
 sudo docker compose build && sudo docker compose up 
 ```

In other terminal window(another environment):  

``` bash
 go test handlers/handlers_test.go && go test dbConn/dbUtils_test.go 
```