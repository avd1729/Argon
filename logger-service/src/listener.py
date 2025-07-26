import pika
import json
from pymongo import MongoClient
import os
from dotenv import load_dotenv

load_dotenv()

def listen():
    RABBIT_URL = os.getenv("RABBIT_MQ_LISTENER_URL")
    MONGODB_CONNECTION_STRING = os.getenv("MONGODB_CONNECTION_STRING")
    MONGO_DB = "logs"
    MONGO_COLLECTION = "log_entries"
    LOG_QUEUE = "log.queue"

    # Connect to MongoDB
    mongo_client = MongoClient(MONGODB_CONNECTION_STRING)
    mongo_db = mongo_client[MONGO_DB]
    log_collection = mongo_db[MONGO_COLLECTION]

    # Connect to RabbitMQ
    connection = pika.BlockingConnection(pika.URLParameters(RABBIT_URL))
    channel = connection.channel()

    channel.queue_declare(queue=LOG_QUEUE, durable=True)

    print(f" [*] Waiting for logs in '{LOG_QUEUE}'. To exit press CTRL+C")

    def callback(ch, method, properties, body):
        try:
            decoded = body.decode("utf-8")
            log_entry = json.loads(decoded)

            if isinstance(log_entry, str):
                log_entry = {
                    "level": "INFO",
                    "message": log_entry,
                    "context": None
                }

            if not isinstance(log_entry, dict):
                raise ValueError("Log entry must be a dict")

            log_collection.insert_one(log_entry)
            print(" [x] Logged:", log_entry)

        except Exception as e:
            print(" [!] Failed to log entry:", e)
            print(" [!] Message body was:", body)


    channel.basic_consume(queue=LOG_QUEUE, on_message_callback=callback, auto_ack=True)

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print(" [x] Stopping log listener...")
        channel.stop_consuming()
        connection.close()
