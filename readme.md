# Clutch

Clutch is a tool for ingesting data from a variety of sources and having refined control over it.

## Features

The clutch software:
- Revieves any kind of data (JSON right now)
- Ingests it (Uses websocket)
- Parses it (JSON right now just stores a map[string]interface{})
- Stores it (Elasticsearch)
<!-- - Correlates it -->
<!-- - Visualizes it -->
- Allows for:
    - x - Masking 
    - x - Synthesizing
    - Anonymizing (Maybe use something like the faking package but I think LLM)
    - Monitoring 
    - Alerting
    - ML on top

- Future:
    - Simulations:
        - Given a situation (i.e. anamolous weather occured in a 30d period which included strong winds and heavy consistent rain, etc)
        - Generate a synthetic dataset which emulates this situation and runs alerting mechanisms on top of it
        - Describe the impact the simulated data shows

To Do Next:
- x - Add RAG similar to here (https://github.com/hantmac/langchaingo-ollama-rag/blob/main/rag/ollama.go)
- x - Build out the model file and stuff for use outside of the storage
- Refine Qdrant cod and the interface for DBs
- Use LLM to build out simulations based on the data
- Add a service to handle the RAG
- Build tests
- Refactor receiver to use interfaces
- Refactor common to use interfaces


## Defining the config

This is the config file that defines the services to be used and the database to be used.

```yaml
server:
  host: "localhost"
  port: "8080"

database:
  type: "elastic"
  cert_location: "http_ca.crt"
  host: "localhost"
  port: "9200"
  user: "username"
  password: "password"

services:
  - storage
    -elastic
    -qdrant (vectors for RAG)
  - masking
  - mask_storage

```

## Starting Ollama 3.2

https://github.com/ollama/ollama?tab=readme-ov-file

```bash 
ollama run llama3.2
```

## Starting QDrant Service

```bash
docker run -p 6333:6333 -p 6334:6334 -v $(pwd)/qdrant_storage:/qdrant/storage:z qdrant/qdrant
```

# Access
REST API: localhost:6333
Web UI: localhost:6333/dashboard
GRPC API: localhost:6334


# Links
https://qdrant.tech/documentation/quickstart/

## Starting the Elasticsearch
```bash
docker run -d --name elasticsearch --net elastic -p 9200:9200 -e "discovery.type=single-node" elasticsearch:8.9.0
```

# Links
https://www.elastic.co/guide/en/elasticsearch/reference/8.15/docker.html
https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/connecting.html
https://docs.go-blueprint.dev/blueprint-core/db-drivers/
https://go.dev/doc/tutorial/database-access
