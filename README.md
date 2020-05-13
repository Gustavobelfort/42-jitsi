# 42-Jitsi (Temporary Name)

![goreportcard](https://goreportcard.com/badge/github.com/Gustavobelfort/42-jitsi)

> Made by [gus](https://github.com/Gustavobelfort) and [pistache](https://github.com/clafoutis42) for 42-Network.

## Introduction

**42-Jitsi** is a golang powered project meant to provide a video conference room to students using [Jitsi](https://jitsi.org) for their remote evaluations.

The repository will be using:
- [gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/jinzhu/gorm)
- [cron.v3](https://github.com/robfig/cron)
- [slack-that](https://github.com/jgengo/slack-that)

> This project was made to facilitate remote evaluations during the COVID-19 crisis.  
> As a distancing measure, we suggest you use it for a while a short time after the re-opening of your campus.

## How it works

![diagram](/assets/diagram.png)

 - After started 42 Jitsi receives information from one of the configured consumers ( see [Consumers](###The-consumers) ) 
 - Stores the processed scale teams into a PostgreSQL DB
 - A daemon runs on a configured wait interval ( 15 minutes by default )
 - At each run the daemon:  
    - Gets the notifiable scale teams from the database
    - Send a HTTP request to a configured Slack_That Client
    - Updates the db to set the scale teams as notified if everything occurs sucessfully

## Usage

### Slack Configuration

You will need to deploy [slack_that](https://github.com/jgengo/slack_that) in your infrastructure.

The bot token will require the scopes `chat:write`, `im:write`, `users:read`, `users:read.email`.

### Configuration

Read the configuration samples _[configs.sample.yaml](./configs/configs.sample.yml)_ and _[example.env](./configs/example.env)_ to understand better
what can be configured.

### Deployment

> In order not to get confused with the multiple docker-compose files, it is advised to use the [script](./docker-compose.sh) `docker-compose.sh`
> to run docker-compose commands. It will use the right docker-compose files depending on the env vars set.  
> Note that this script will try to read the `.env` file.

**Requirements:**
- `docker`
- `docker-compose`

#### Production

You need to set the env var `ENVIRONMEMT` to `production`. e.g:
```
staff@42campus:~/42-jitsi # ENVIRONMENT=production ./docker-compose.sh up -d
Creating network "42-jitsi_42jitsi" with the default driver
Creating 42-jitsi_consumer_1_22f48fd4f03b ... done
staff@42campus:~/42-jitsi # ENVIRONMENT=production ./docker-compose.sh ps
              Name                 Command    State         Ports
-------------------------------------------------------------------------
42-jitsi_consumer_1_2aca4c53e112   /bin/api   Up      127.0.0.1->5000/tcp
staff@42campus:~/42-jitsi # ENVIRONMENT=production ./docker-compose.sh down
Removing 42-jitsi_consumer_1_2aca4c53e112 ... done
Removing network 42-jitsi_42jitsi
```

This will deploy the server ready to use in production.  
The [api](./cmd/api) consumer is deployed by default.

#### Development

This is the default behaviour. You need to set env var `ENVIRONMENT` to anything else than `production`. e.g:
```
staff@42campus:~/42-jitsi # ENVIRONMENT=local ./docker-compose.sh up -d
Creating network "42-jitsi_42jitsi" with the default driver
Creating 42-jitsi_rabbitmq_1_724cafefbccd ... done
Creating 42-jitsi_consumer_1_c8521bca9d9f ... done
Creating 42-jitsi_db_1_24d49a369d71       ... done
staff@42campus:~/42-jitsi # ENVIRONMENT=local ./docker-compose.sh ps
              Name                            Command               State                      Ports
-----------------------------------------------------------------------------------------------------------------------
42-jitsi_consumer_1_496f749bef9f   /bin/api                         Up       127.0.0.1:5432->5432/tcp
42-jitsi_db_1_9120aefbab27         docker-entrypoint.sh postgres    Up       127.0.0.1:5000->5000/tcp
42-jitsi_rabbitmq_1_ade8110994b6   docker-entrypoint.sh rabbi ...   Up       15671/tcp, 127.0.0.1:15672->15672/tcp,
                                                                             25672/tcp, 4369/tcp, 5671/tcp,
                                                                             127.0.0.1:5672->5672/tcp
staff@42campus:~/42-jitsi # ./docker-compose.sh down
Stopping 42-jitsi_rabbitmq_1_ade8110994b6 ... done
Removing 42-jitsi_db_1_9120aefbab27       ... done
Removing 42-jitsi_consumer_1_496f749bef9f ... done
Removing 42-jitsi_rabbitmq_1_ade8110994b6 ... done
Removing network 42-jitsi_42jitsi
```
It will deploy the service as standalone. It will have its own postgresql and rabbitmq container and force set
the corresponding environmental variables so that your container connects to them.

### The consumers

There are different type of consumers that you can use. Here's an exhaustive list of them:
- [api](./cmd/api)
- [rabbit](./cmd/rabbit)

You can choose which one to deploy with docker-compose by setting the env var `JISTI42_CONSUMER_TYPE` to the corresponding
value.

_Common configurations_:
- General configuration: `TIMEOUT`
- Logging configuration: `LOG_LEVEL`, `SENTRY_DSN`, `SENTRY_LEVELS`, `SENTRY_ENABLED`, `LOGSTASH_HOST`, `LOGSTASH_PORT`
  `LOGSTASH_PROTOCOL`, `LOGSTASH_LEVELS`, `LOGSTASH_ENABLED`
- Intranet application: `INTRA_APP_ID`, `INTRA_APP_SECRET`.
- PostgreSQL database: `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`.

#### API Consumer

> See more details [here](./cmd/api).

To deploy it with docker-compose, set `JITSI42_CONSUMER_TYPE` to `api`.

This consumer will expose an api on the port `5000` _(by default)_ with [gin](https://github.com/gin-gonic/gin) that will
expect payloads from **42's intranet scale_team webhooks**.

_Specific configurations:_
- Exposure address: `HTTP_ADDR`.
- Intranet webhooks: `INTRA_WEBHOOKS`.

#### Rabbit Consumer

> See more details [here](./cmd/rabbit).

To deploy it with docker-compose, set `JITSI42_CONSUMER_TYPE` to `rabbit`.

This consumer will read from a [rabbitmq](https://www.rabbitmq.com/) queue. The messages' bodies are expected to be payloads
from **42's intranet scale_team webhooks**. Of course the corresponding headers are expected to be set.

_Specific configurations:_
- RabbitMQ: `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_VHOST`, `RABBITMQ_USER`, `RABBITMQ_PASSWORD`, `RABBITMQ_QUEUE`.