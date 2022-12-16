# amqpresurrector

Auto Reconnect rabbitmq connection. It's wrapper so most methods should work.

After reconnection, these operations will be exucted again, on failure all reconnection mechnism
will fail:

- QueueDeclare
- ExchangeDeclare
- Restore Channels
- Consume delivery channel doesn't close on disconnection
- QueueBind

please check [example](./example).

## Disclaimer

it's not stable version, DON'T RECOMMEND TO USE IN PRODUCTION.

## LICENSE
