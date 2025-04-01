# Distributed FileSystem

    server will recieve file related request from clients and simultaneously updated across all distributions

## Server Creation

    1. a listener is created for a given address
    2. transport will continue to listen for new connections
    3. a recieved connection is handled to recieve incoming messages
    4. each message recieved on a conn is sent to server via channel
    5. data is encoded is stored in filesystem and broadcasted to bootstraped nodes

## Handling Files

    gob encoding
