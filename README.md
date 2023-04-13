
To test, run

`docker-compose up -d`

The apis for adding and retrieving messages are different services, and thus on different ports. The api for sending is exposed on `localhost:8080`, and the api for retreiving is on `localhost:8081`.


The following commands can be executed to test the functionality:

```
curl --request POST \
  --url http://localhost:8080/message \
  --header 'Content-Type: application/json' \
  --data '{
	"message": "first message",
	"sender": "karl",
	"receiver": "maria"
}'

curl --request POST \
  --url http://localhost:8080/message \
  --header 'Content-Type: application/json' \
  --data '{
	"message": "second message",
	"sender": "karl",
	"receiver": "maria"
}'

curl --request GET \
  --url 'http://localhost:8081/message/list?sender=karl&receiver=maria'

```
  
Things that I would like to add but did not have time includes:
- unit tests
- integration tests
- automatic reconnect on connection loss
- structured logging
- separate transport layer of APIs from main modules
- better error handling and verbose http responses/statuses
