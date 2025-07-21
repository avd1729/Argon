package com.example.webhooklistener.queue.listener;

import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;

@Component
public class RabbitMqListener {

    @RabbitListener(queues = "webhook.queue")
    public void handleMessage(String message) {
        System.out.println("Processing webhook from queue: " + message);
    }
}

