# CMPE-273-Assignment2

### POST Curl command

curl -H "Content-Type: application/json" -X POST -d '{"name" : "juan", "address" : "123 Main St", "city" : "San Francisco","state" : "CA","zip" : "94113"}' localhost:8080/locations

### Output
{"ID":"562c7a940ed561698473bacb","name":"juan","address":"123 Main St","city":"San Francisco","state":"CA","zip":"94113","coordinate":{"Lat":37.7917618,"Lng":-122.3943405}}


### GET curl command

curl -i localhost:8080/locations/562c7a940ed561698473bacb

### Output

{"ID":"562c7a940ed561698473bacb","name":"juan","address":"123 Main St","city":"San Francisco","state":"CA","zip":"94113","coordinate":{"Lat":37.7917618,"Lng":-122.3943405}}

### PUT curl command

curl -H "Content-Type: application/json" -X PUT -d '{"address" : "123 Main St", "city" : "San Francisco","state" : "CA","zip" : "94113"}' localhost:8080/locations/562c7a940ed561698473bacb

### Output

{"ID":"562c7a940ed561698473bacb","name":"juan","address":"123 Main St","city":"San Francisco","state":"CA","zip":"94113","coordinate":{"Lat":37.7917618,"Lng":-122.3943405}}

### DELETE curl command

curl -H "Content-Type: application/json" -X DELETE localhost:8080/locations/562c7a940ed561698473bacb

### Output

Record Deleted.





