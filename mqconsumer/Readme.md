## RabbitMQ 

### Priority
For Priority Queue we have to add 
`"x-max-priority" : 10` in our queue to make it support but 
priority works in a different way when a new consumer starts 
to consume then it gets high priority message first but if a 
Consumer is working on the queue the sequences is maintained 
(ignoring priority) this poc is with direct exchange


###  Message Timeout
https://www.rabbitmq.com/ttl.html#per-message-ttl-caveats
http://www.inanzzz.com/index.php/post/1p7m/creating-a-rabbitmq-dlx-dead-letter-exchange-example-with-golang
https://www.rabbitmq.com/dlx.html#overview
