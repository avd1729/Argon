package com.example.webhooklistener.queue.listener;

import com.example.webhooklistener.webhook.dto.WebhookPayload;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;

@Component
public class RabbitMqListener {

    @RabbitListener(queues = "webhook.queue")
    public void handleMessage(WebhookPayload payload) {
        System.out.println("Received webhook for repo: " + payload.getRepoName() + ", branch: " + payload.getBranch());
    }
}
