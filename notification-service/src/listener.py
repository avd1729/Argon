import pika
from src.email_service import send_email

def listen(queue_name='notification.queue', host='localhost'):
    def callback(ch, method, properties, body):
        send_email(body.decode())

    # Connect to RabbitMQ
    connection = pika.BlockingConnection(pika.ConnectionParameters(host=host))
    channel = connection.channel()

    # Ensure the queue exists
    channel.queue_declare(queue=queue_name)

    # Start consuming messages
    channel.basic_consume(
        queue=queue_name,
        on_message_callback=callback,
        auto_ack=True
    )

    print(f"[*] Listening on queue '{queue_name}'. Press CTRL+C to stop.")
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print("\n[!] Stopped by user")
    finally:
        channel.stop_consuming()
        connection.close()
