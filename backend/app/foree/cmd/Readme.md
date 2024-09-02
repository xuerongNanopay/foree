run local: `go build -o main && ./main app foree_app`
test:
`curl http://localhost:8080/sys/v1/hello`
`curl http://localhost:8080/sys/v1/mysql_connection`
run local with custome env-file: `go build -o main && ./main app foree_app --env-file  ../deploy/.local_env`