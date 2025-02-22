# Receipt API

This is a simple webservice that fulfills a receipt API. There are two supported endpoints.

### Endpoint: Process Receipts
* Path: '/receipts/process'
* Method: 'POST'
* Payload: Receipt JSON
* Response: JSON containing an id for the receipt.

Description:

Takes in a JSON receipt and returns a JSON object with a generated unique ID.

### Endpoint: Get Points
* Path: '/receipts/{id}/process'
* Method: 'GET'
* Response: JSON containing number of points for the receipt.

Description:

Looks up receipt by the ID and returns an object specifying points awarded following specified rules.

## Instructions to run

There are two ways to run this webservice. The first is with Docker. Make sure Docker is running so it can connect. In the project directory, to run "docker build --tag docker-receipt-api ." to build a docker image.  Then run with "docker run -p 8000:8000 docker-receipt-api". The api will then be running on port 8000, to which you can send the support GET and POST requests.

Alternatively, you can run this application with "go run main.go". You may need to get a couple of things beforehand: "go get github.com/google/uuid" and "go get github.com/gorilla/mux".
