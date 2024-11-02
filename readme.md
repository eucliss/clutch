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
    - Masking
    - Synthesizing
    - Anonymizing (Faking package)
    - Monitoring
    - Alerting
    - ML on top

- Future:
    - Simulations:
        - Given a situation (i.e. anamolous weather occured in a 30d period which included strong winds and heavy consistent rain, etc)
        - Generate a synthetic dataset which emulates this situation and runs alerting mechanisms on top of it
        - Describe the impact the simulated data shows