# API caching middleware server
Redis middleware cache for Autopilot contact GET/POST API. 

# Download and build the project  

## Downloading 
Downloading the project:

`go get github.com/romangurevitch/redis-cache-go`

## Testing the project:
To test the project run:

`go test github.com/romangurevitch/redis-cache-go/...`

## Running the project:
In order to run the server redis instance must run and configured correctly in the config file. 

To run the contact caching server: 

`go run contact/contact.go`

## Config file 
To change redis, server or API config settings modify values in `config.go` file.  