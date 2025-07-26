package com.example.webhooklistener.queue.producer;

import com.example.webhooklistener.enums.RMQueue;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Service;

@Service
public class Producer {

    private final RabbitTemplate rabbitTemplate;

    public Producer(RabbitTemplate rabbitTemplate) {
        this.rabbitTemplate = rabbitTemplate;
    }

    public void pushToQueue(RMQueue queue, Object payload){
        String queueName = queue.getQueueName();
        rabbitTemplate.convertAndSend(queueName, payload);
    }
}
