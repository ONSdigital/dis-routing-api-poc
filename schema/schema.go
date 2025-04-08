package schema

import "github.com/ONSdigital/dp-kafka/v3/avro"

var routingUpdated = `{
  "type": "record",
  "name": "routing-updated",
  "fields": [
    {"name": "event_type", "type": "string", "default": ""},
    {"name": "changes", "type": {"type": "array", "items": {
      "name": "Changes",
      "type" : "record",
      "fields": [
        { "name": "action", "type": "string", "default": "" },
        { "name": "entity", "type": "string", "default": "" },
        { "name": "id", "type": "string", "default": "" },
        { "name": "timestamp", "type": "string", "default": "" }
      ]
    }}, "default": []},
  ]
}`

// RoutingUpdatedEvent is the Avro schema for Routing Update messages.
var RoutingUpdatedEvent = &avro.Schema{
	Definition: routingUpdated,
}
