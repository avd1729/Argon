package com.example.webhooklistener.webhook.service;

import com.example.webhooklistener.enums.RMQueue;
import com.example.webhooklistener.queue.producer.Producer;
import com.example.webhooklistener.webhook.dto.WebhookPayload;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.Map;

@Service
public class WebhookService {

    private final Producer producer;

    public WebhookService(Producer producer) {
        this.producer = producer;
    }

    public void addToQueue(Map<String, Object> payload) {
        try {
            WebhookPayload payload1 = parseGitHubWebhook(payload);
            // Log the payload
            producer.pushToQueue(RMQueue.LOGGER_QUEUE, payload1);
            producer.pushToQueue(RMQueue.WEBHOOK_QUEUE, payload1);
            String response = "Webhook received and pushed to RabbitMQ";
            producer.pushToQueue(RMQueue.LOGGER_QUEUE, response);
        } catch (Exception e) {
            // Log the exception
            producer.pushToQueue(RMQueue.LOGGER_QUEUE, e);
            throw new RuntimeException(e);
        }

    }

    public WebhookPayload parseGitHubWebhook(Map<String, Object> payload) {
        WebhookPayload dto = new WebhookPayload();

        Map<String, Object> repository = (Map<String, Object>) payload.get("repository");
        Map<String, Object> headCommit = (Map<String, Object>) payload.get("head_commit");
        Map<String, Object> pusher = (Map<String, Object>) payload.get("pusher");

        dto.setRepositoryUrl((String) repository.get("clone_url"));
        dto.setRepoName((String) repository.get("name"));
        dto.setBranch(((String) payload.get("ref")).replace("refs/heads/", ""));
        dto.setCommitId((String) headCommit.get("id"));
        dto.setCommitMessage((String) headCommit.get("message"));
        dto.setPusherName((String) pusher.get("name"));
        dto.setPusherEmail((String) pusher.get("email"));
        dto.setTimestamp(Instant.now().toString());
        dto.setTriggerSource("push");
        dto.setProjectType("maven");

        return dto;

    }

}
