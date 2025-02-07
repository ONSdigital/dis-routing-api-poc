# dis-routing-api-poc

Proof of concept for a Router API that would be the database interface to CRUD path based routes and redirects.

## Features

* Full CRUD operations for managing routes and redirects.
* Support for multiple domains.
* Validation to prevent overlapping or circular routes.
* Updates handled using Mongo/DocumentDB change streams, a reload endpoint, or cache update intervals.
* Wildcard support for flexible route and redirect configurations.
* Request logging and analytics for tracking redirects and route hits.
* OpenTelemetry integration for monitoring.

## API Endpoints

### Routes Management

* `POST /routes` - Create a new route.
* `GET /routes/{id}` - Retrieve a specific route.
* `PUT /routes/{id}` - Update an existing route.
* `DELETE /routes/{id}` - Remove a route.

### Redirects Management

* `POST /redirects` - Create a new redirect.
* `GET /redirects/{id}` - Retrieve a specific redirect.
* `PUT /redirects/{id}` - Update an existing redirect.
* `DELETE /redirects/{id}` - Remove a redirect.

### Traffic Routing

* `GET /{path}` - Handle incoming requests:
    * If a matching route exists, proxy the request to the destination URL.
    * If a redirect exists, return an HTTP 301/302 response.
    * If no match is found, return a 404 Not Found response.

### Admin & Management

* `POST /admin/reload` - Reload configuration.
* `GET /admin/health` - Health check endpoint.

## Technology Stack

* Language: Go
* Framework: chi for routing
* Database: MongoDB / DocumentDB (mongo-go-driver for interactions)
* Monitoring: OpenTelemetry for tracing and metrics

## Getting started

* Run `make debug` to run application on http://localhost:29700
* Run `make help` to see full list of make targets

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default             | Description                                                                                                        
|------------------------------|---------------------|--------------------------------------------------------------------------------------------------------------------
| BIND_ADDR                    | :29700              | The host and port to bind to                                                                                       
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                  | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  
| HEALTHCHECK_INTERVAL         | 30s                 | Time between self-healthchecks (`time.Duration` format)                                                            
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                 | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) 
| OTEL_EXPORTER_OTLP_ENDPOINT  | localhost:4317      | Endpoint for OpenTelemetry service                                                                                 
| OTEL_SERVICE_NAME            | dis-routing-api-poc | Label of service for OpenTelemetry service                                                                         
| OTEL_BATCH_TIMEOUT           | 5s                  | Timeout for OpenTelemetry                                                                                          
| OTEL_ENABLED                 | false               | Feature flag to enable OpenTelemetry                                                                               

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2025, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
