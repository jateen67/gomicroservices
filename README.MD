# Distributed System

![0-tourcropped](https://github.com/jateen67/mlbstatsapi/assets/106696411/977db707-20b1-4ecb-a0b2-bd1042f838b5)

The purpose of this project was to build a scalable application using microservice architecture.

### What is microservice architecture?

Microservices, often known as microservices architecture, is a method of building large applications out of smaller, modular components or services. A given service can handle functionality like user identification/authentication, alert generation, data logging, and a number of other functions. Each service handles its own database and typically runs a distinct process.

### What are some advantages of using microservice architecture?

Microservices have a wide variety of benefits:

- They can each be developed and deployed individually without having to interfere with/affect any other part of the application, unlike monolithic applications

- They can be written using different languages and tools and still communicate with one another. They are very flexible in that sense

- They have the ability to scale quickly and be reused across different projects

- They lead to better and more efficient fault isolation and resiliency

- They can be owned and maintained by a relatively small team of engineers

- They work well with containers (Docker, etc.)

### How to run

1. `cd` into the `project` directory
2. Run `docker-compose up -d`
3. `cd` into the `client` directory
4. Run `npm install`, followed by `npm run dev`

### Services

[architecture picture]

**Front-end Service**

The most simple service. It's simply a web page created in React.js+TypeScript (and Bootstrap for styling) that allows the user to interact with the various services.

The user can communicate with the services using various methods, such as via REST, gRPC protocol, sending/consuming messages using AMQP (Advanced Message Queuing Protocol), and also SMTP (Simple Mail Transfer Protocol).

**Broker Service**

This serves as a main point of entry into the microservice cluster. It is optional, but in this application, all requests to the various microservices from the client go through the Broker service first, which then communicates with the respective services.

When a user sends a request to a certain microservice, the request goes to the Broker service first, which in turn will communicate and send its own request to that specific service, and finally return the response to the client.

Since the Broker service's main purpose is to communicate with the other services, sending a request to it directly will give a basic response back to the user, as a simple way of indicating that it is functioning well and ready to communicate.

![1-brokercropped](https://github.com/jateen67/mlbstatsapi/assets/106696411/8fe8669a-2b40-4e36-9c03-065e3a5a8279)

**Logger Service**

This is a service that logs some kind of simulated activity, whatever that may be. When the user sends a request, it will simply insert some data into a MongoDB database, indicated that the user has done some activity and it has been successfully stored/logged, kind of like a traditional activity logger in any other application.

When the user sends a request to the Logger Service through the Broker Service, the two services communicate with one other via gRPC.

The database containing the successful user activity logs can be accessed locally using a Mongo client like [MongoDBCompass](https://www.mongodb.com/products/compass) (Connection String: mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&directConnection=true&ssl=false)

![3-grpccropped](https://github.com/jateen67/mlbstatsapi/assets/106696411/cf077641-8fa7-47e0-b4f2-92a8baf4ec1e)

**Authentication Service**

This is a service that simulates attempting to "log in" a user given the proper credentials. When the user sends a request, a username and password is sent alongside the request (admin@example.com/verysecret). The service will take these credentials and compare them to the hashed credentials it has stored in a Postgres database to try to find a match. If there is a match, then a success message will be sent to the client.

The Authentication Service will communicate with the Logger Service to store that successful user login attempts, akin to what you might see in a traditional activity logger in any other application.

When the user sends a request to the Authentication Service through the Broker Service, the two services communicate with one other through JSON to fulfill the user request via a RESTful API.

The database containing the user credentials can be accessed locally using a lightweight database manager like [Beekeeper Studio](https://www.beekeeperstudio.io/) (Connection String: host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5)

The database containing the successful user login logs, alongside the other successful activity logs, can be accessed locally using a Mongo client like [MongoDBCompass](https://www.mongodb.com/products/compass) (Connection String: mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&directConnection=true&ssl=false)

![2-authcropped](https://github.com/jateen67/mlbstatsapi/assets/106696411/60fc90e4-8f7b-4ae0-bad3-33bc8b56724a)

**Listener Service**

This service is another way that the user can send a request to a Logger Service to store a simulated activity. It accomplishes the exact same things as the standard Logger Service, but through a different method.

When the user sends a request to the Logger Service to log some simulated activity via this alternative method, the Broker Service will not communicate with the Logger Service directly like normal, but will instead communicate with the Listener Service, which will then in turn communicate to the Logger Service through RabbitMQ.

This works by first pushing an event to the RabbitMQ server from the Broker Service via AMQP. RabbitMQ then takes that event and adds it to a queue. The Listener Service looks at that queue and constantly checks to see if there are any messages it should read. If so, it reads it, decides what it needs to do with it depending on the content it reads, and then calls the appropriate service to perform the action. In this case, the Logger Service gets called to store the simulated activity it received from the message into its MongoDB database.

The database containing the successful user activity logs via RabbitMQ can be accessed locally using a Mongo client like [MongoDBCompass](https://www.mongodb.com/products/compass) (Connection String: mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&directConnection=true&ssl=false)

![4-rabbitcropped](https://github.com/jateen67/mlbstatsapi/assets/106696411/a62e37fa-cfb0-4cdb-a3a8-31233f955b1b)

**Mail Service**

This is an email service that allows the user to send an email out with MailHog. It works by taking in a JSON payload from the client, converting it into a formatted email, then sending it via SMTP.

Sent emails can be viewed by accessing the MailHog UI on `localhost:8025`

![5-mailercropped](https://github.com/jateen67/mlbstatsapi/assets/106696411/4dfa002a-b609-436f-bdf1-8b3d5ff125ed)

### Technologies

These were the main tools used to build this project:

- [React.js](https://react.dev/)
- [TypeScript](https://www.typescriptlang.org/)
- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [MongoDB](https://www.mongodb.com/)
- [gRPC](https://grpc.io/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [Bootstrap](https://getbootstrap.com/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
