# dis-routing-api-poc

---

:warning: This repository was archived in May 2025 and is no longer in development. :warning:

---

Proof of concept for a Router API that would be the database interface to CRUD path based routes and redirects.

Please refer the to documentation here - https://confluence.ons.gov.uk/display/DIS/Routing+APIs for Objectives and
design considerations.

### Technology Stack

* Language: Go
* Framework: mux for routing
* Database: MongoDB / DocumentDB (mongo-go-driver for interactions)**

## Getting started

* Run `make debug` to run application on http://localhost:29700
* Run `make help` to see full list of make targets

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                                               | Description                                                                                                        
|------------------------------|-------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------
| BIND_ADDR                    | :29700                                                | The host and port to bind to                                                                                       
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                                                    | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  
| HEALTHCHECK_INTERVAL         | 30s                                                   | Time between self-healthchecks (`time.Duration` format)                                                            
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                                                   | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) 
| MONGODB_BIND_ADDR            | localhost:27017                                       | The MongoDB bind address                                                                                           
| MONGODB_USERNAME             |                                                       | MongoDB Username                                                                                                   
| MONGODB_PASSWORD             |                                                       | MongoDB Password                                                                                                   
| MONGODB_DATABASE             | router                                                | The MongoDB router database                                                                                        
| MONGODB_COLLECTIONS          | RoutesCollection:routes,RedirectsCollection:redirects | MongoDB collections                                                                                                
| MONGODB_ENABLE_READ_CONCERN  | false                                                 | Switch to use (or not) majority read concern                                                                       
| MONGODB_ENABLE_WRITE_CONCERN | true                                                  | Switch to use (or not) majority write concern                                                                      
| MONGODB_CONNECT_TIMEOUT      | 5s                                                    | The timeout when connecting to MongoDB (`time.Duration` format)                                                    
| MONGODB_QUERY_TIMEOUT        | 15s                                                   | The timeout for querying MongoDB (`time.Duration` format)                                                          
| MONGODB_IS_SSL               | false                                                 | Switch to use (or not) TLS when connecting to mongodb                                                              

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2025, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
