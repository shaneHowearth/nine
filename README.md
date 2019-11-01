# Article API

# Setup/Installation
Require `docker` and `docker-compose`
The easiest way to get things up and running is `docker-compose up`

# Notes
This is a RESTful API that provides the ability to Create (but not Update or Delete) Articles.

To make it interesting I have created a CQRS system, with a Command path where articles are created, and persisted, and a Query path where articles or info about them are requested, either by using tags, or ids as parameters.

The persistence for this implementation is achieved using a PostgreSQL database (Note: The docker container for the database is NOT setup to persist data permanently on the host/test system. The stopping or removing of the PostgreSQL container will result in loss of all data).
There are no indexes on the database, to speed up inserts. This slows down searches, but this is helped by the fact that the searches are performed in the NoSQL servers.

The business layer for each microservice defines a repository interface that determines what the business logic will communicate to the repository, allowing any repository software to be used, as long as the interface is implemented. This means that the business logic never need know if the repository layer is changed, but the repository layer MUST be informed if the way the business layer wants to communicate to the repository is changed (Dependency Inversion Principle).

Because the Data in a CQRS system is, by definition, Partitioned, I have a choice to make whether to make the data consistent, or available. I have chosen consistency for getting articles by ID, and availability for tags. This means that getting an article by ID is slow, but guaranteed to be authoritative. If an article for a given ID is not found, that is because, at the time of making the request, it does not exist. Tags search is faster, but not guaranteed to be consistent, it's "eventually consistent".

The database that holds the articles is in its own container, this allows multiple createarticle processes to exist at the same time, each writing to the db (the database will handle the multithreading issues). But the key reason is that the database container will have code that depends on the read containers. That is the communication will be defined in the code bases that the read containers applications are built from. Seperation to another container means that the create containers will not be burdened with the dependency.

# Assumptions
There are no updates or deletes of Articles. Once an Article is created, only another can be added.
The author of the article hasn't been recorded, if that were to happen then authentication and other author related data might need to be collected, in a user system.
There is NO authentication nor authorisation specified, adding such functionality would not be trivial.
There is NO encryption between services.
I have limited tag names to be characters of the set [0-9a-zA-Z ].

# Integration

Communication between services is gRPC and RabbitMQ for this project. As a Remote Procedure invocation, gRPC is fast, and can be allowed to cross networks (firewalled or other), however it should be noted that it couples services together, if the server changes the signature of the method called, then all the clients need to be updated, this is mitigated by defining the client with the server code, and importing the definition, but that creates a code dependency. Also there is possibility for the call to be lost inflight, with neither the server, nor the client, being aware of the failure. This is mitigated by establishing a timeout on the client side, and a retry, however other options exist, such as circuit breakers, and fallbacks for situations where the server is unable to respond for a variety of reasons.
RabbitMQ is used as a Pub/Sub, but this can be changed to a fan out if multiple caches are present.

# What did I think of the test

On the surface it's a simple REST server with a few verbs implemented, however in order to make a scalable system a lot of work needs to be done to ensure that the integration between services is robust, because there's no guarantee that the other service is up, available, or doing the job it's expected to do, or even that the network is behaving properly. I'm enjoying it thoroughly =)
