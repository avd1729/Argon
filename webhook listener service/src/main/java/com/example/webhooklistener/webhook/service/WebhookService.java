package com.example.webhooklistener.webhook.service;

import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Service;

@Service
public class WebhookService {
    private final RabbitTemplate rabbitTemplate;

    public WebhookService(RabbitTemplate rabbitTemplate) {
        this.rabbitTemplate = rabbitTemplate;
    }
    public void addToQueue(String payload) {
        rabbitTemplate.convertAndSend("webhook.queue", payload);
    }
}
