# Distributed FileSystem

    server will recieve file related request from clients and simultaneously updated across all distributed servers

## Server Creation

    1. a listener is created for a given address
    2. transport will continue to listen for new connections
    3. a recieved connection is handled to recieve incoming messages
    4. each message recieved on a conn is sent to server via channel
    5. data is stored in filesystem and broadcasted to bootstraped nodes
    6. data transfer to bootstrap nodes is through gob enconding
    7. single leader multiple replica model
