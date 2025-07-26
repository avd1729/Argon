import pika
from src.email_service import send_email

import os
from dotenv import load_dotenv

load_dotenv()

RABBIT_URL = os.getenv("RABBIT_MQ_LISTENER_URL")
NOTIFICATION_QUEUE = "notification.queue"


def listen():
    def callback(ch, method, properties, body):
        try:
            send_email(body.decode())
            print("[x] Email sent:", body.decode())
        except Exception as e:
            print(f"[!] Error processing message: {e}")

    def connect_and_consume():
        connection = pika.BlockingConnection(pika.URLParameters(RABBIT_URL))
        channel = connection.channel()
        channel.queue_declare(queue=NOTIFICATION_QUEUE, durable=True)

        channel.basic_consume(
            queue=NOTIFICATION_QUEUE,
            on_message_callback=callback,
            auto_ack=True
        )

        print(f"[*] Listening on queue '{NOTIFICATION_QUEUE}'. Press CTRL+C to stop.")
        channel.start_consuming()

    while True:
        try:
            connect_and_consume()
        except KeyboardInterrupt:
            print("\n[!] Stopped by user")
            break
        except Exception as e:
            print(f"[!] Listener crashed: {e}. Reconnecting...")
