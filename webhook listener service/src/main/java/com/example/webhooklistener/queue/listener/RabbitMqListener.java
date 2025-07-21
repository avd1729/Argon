package com.example.webhooklistener.queue.listener;

import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;

@Component
public class RabbitMqListener {

    @RabbitListener(queues = "test.queue")
    public void listen(String message) {
        System.out.println("Received from RabbitMQ: " + message);
    }
}

