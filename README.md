# API caching middleware server
Redis middleware cache for Autopilot contact GET/POST API. 

# Download and build the project  

## Downloading 
Downloading the project:

`go get github.com/romangurevitch/redis-cache-go`

## Testing the project
To test the project run:

`go test github.com/romangurevitch/redis-cache-go/...`

## Running the project
In order to run the server redis instance must run and configured correctly in the config file. 

To run the contact caching server: 

`go run contact/contact.go`

## Config file 
To change redis, server or API config settings modify values in `config.go` file.  

## Usage examples
In order to run these examples Autopilot API key is required. 

Full API documentation can be found [here](https://autopilot.docs.apiary.io/).

### Create a contact
Create a `contact.json` file with the following content: 

```json
{
  "contact": {
    "FirstName": "Slarty",
    "LastName": "Bartfast",
    "Email": "test@slarty.com",
    "custom": {
      "string--Test--Field": "This is a test"
    }
  }
}
```

Run the following command, replace `APIKEY` header value.  

```shell script
curl -H "Content-Type: application/json" -H "autopilotapikey: <APIKEY>" -XPOST -d @contact.json http://localhost:8080/contact
```

### Getting a contact
Run the following command, replace `APIKEY` header value.  

```shell script
curl -H "autopilotapikey: <APIKEY>" http://localhost:8080/contact/test@slarty.com
```