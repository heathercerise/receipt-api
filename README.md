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
