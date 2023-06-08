# [☰](../README.md) Instance

In instance block you can check information about particular instance. It mostly works
with postgres itself, sending some request to postgres to monitor and troubleshoot it.

### Chart

In __Chart__ tab you can see some simple charts with information like database uptime, number of 
databases, connections and it sizes.

![instance_chart](images/instance_chart.png)

### Queries

In __Queries__ tab you can run some queries that will help you to troubleshoot postgres. If you have
configured _Manual queries_ when you started postgres, you will be able to create your own queries and
edit system queries to fit your needs. But it can be dangerous, because Ivory executes any query like 
UPDATE / INSERT / DELETE and you will be able to fetch any data. That is why sometimes it is good to restrict
such usage (by default it is configured as disabled)

![instance_queries](images/instance_queries.png)
