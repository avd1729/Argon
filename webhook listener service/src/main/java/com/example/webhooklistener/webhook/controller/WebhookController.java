package com.example.webhooklistener.webhook.controller;

import com.example.webhooklistener.webhook.service.WebhookService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class WebhookController {
    private final WebhookService webhookService;

    public WebhookController(WebhookService webhookService) {
        this.webhookService = webhookService;
    }

    @PostMapping("/api/webhook")
    public ResponseEntity<String> handleWebhook(@RequestBody String payload) {
        System.out.println("Webhook received: " + payload);
        webhookService.addToQueue(payload);
        return ResponseEntity.ok("Webhook received and pushed to RabbitMQ");
    }
}
