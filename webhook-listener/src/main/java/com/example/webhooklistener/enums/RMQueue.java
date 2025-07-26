package com.example.webhooklistener.enums;

public enum RMQueue {
    WEBHOOK_QUEUE("webhook.queue"),
    LOGGER_QUEUE("log.queue");

    private final String queueName;

    RMQueue(String queueName) {
        this.queueName = queueName;
    }

    public String getQueueName() {
        return queueName;
    }
}
