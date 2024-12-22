# LMC


go run main.go

curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
"expression": "2+2*2"
}'


curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
"expression": "2+2//2"
}'
