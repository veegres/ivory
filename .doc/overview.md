# [☰](../README.md) Overview

Overview tab helps you to monitor your cluster - check basic info, statuses, etc. It
provides simple UI for patroni, so you don't need to hesitate with api or cli tool. You
can do such things as:
- switchover (schedule switchover)
- failover
- reinitialisation
- restart (schedule restart)
- reload

![overview](images/overview.png)

You can dig deeply into problematic instance just by simply clicking on it. There you will
find [Query Builder](instance.md) with requests to troubleshoot your database.

### Cluster Options

In the top right corner you can find icons that help to understand what configurations
for this particular cluster you have set up.

You can change setup by clicking on the settings button. There you can specify some options
like set up postgres password or patroni password, add tags or configure mutual connection
between Ivory and patroni or postgres.

Main instance is a specific instance to which Ivory sends request. Automatically
it is always try to choose master as a main instance. If you want you can always change it
here.

![options](images/cluster_options.png)
