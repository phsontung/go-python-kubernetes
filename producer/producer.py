import pika
import random
import sys
import time

import logging

LOGGER = logging.getLogger(__name__)

routing_keys = [
    'kern.*',
    'app.critical',
    '*.info'
]
messages = [
    'Hello world !',
    'Today is Friday',
    'An empty message from producer'
]

LOGGER.warning('Prepare to start producer .....')
time.sleep(20) # Hack here to wait for rabbitmq when run with docker-compose
# connection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq'))
connection = pika.BlockingConnection(pika.URLParameters('amqp://guest:guest@rabbitmq:5672/'))
channel = connection.channel()

channel.exchange_declare(exchange='topic_logs', exchange_type='topic')

while True:
    r_key = random.randint(0, 2)
    m_key = random.randint(0, 2)
    channel.basic_publish(
        exchange='topic_logs', routing_key=routing_keys[r_key], body=messages[m_key])
    LOGGER.warning(" [x] Sent %r:%r" % (routing_keys[r_key], messages[m_key]))
    time.sleep(5)

connection.close()