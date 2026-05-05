cmd folder has 
>api (contains rest endpoints)
>scheduler (used to store jobs in a particular queue)
>worker (these nodes comsume the queue and process the requests in async manner)
all these can be invoked and they run on go run command to execute.