Extension for logrus to enable logging to amqp (rabbit)

Based on https://github.com/cloudaccessio/logrus


For testing first export the environment variable:
TEST_CONNECTION="amqp://<USER>:<PASSWORD>@<HOST>:<PORT>/<OPTIONAL_VHOST>"


example
TEST_CONNECTION="amqp://logg_user:log_super_secret_password@some_host:5672/vhost1"


-- leffen
