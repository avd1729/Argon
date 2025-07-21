package com.example.webhooklistener.webhook.service;

import com.example.webhooklistener.webhook.dto.WebhookPayload;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.Map;

@Service
public class WebhookService {

    private final RabbitTemplate rabbitTemplate;

    public WebhookService(RabbitTemplate rabbitTemplate) {
        this.rabbitTemplate = rabbitTemplate;
    }

    public void addToQueue(Map<String, Object> payload) {
        WebhookPayload payload1 = parseGitHubWebhook(payload);
        rabbitTemplate.convertAndSend("webhook.queue", payload1);
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
