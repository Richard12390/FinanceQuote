package com.org.quoteservice.shared.ws;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@RequiredArgsConstructor
public class QuoteStreamListener {

    private final SimpMessagingTemplate messagingTemplate;
    private final ObjectMapper objectMapper;

    @RabbitListener(queues = "#{quotesStreamQueue.name}")
    public void onMessage(Message message) {
        try {
            QuoteMessage quote = objectMapper.readValue(message.getBody(), QuoteMessage.class);
            messagingTemplate.convertAndSend("/topic/quotes", quote);
        } catch (Exception ex) {
            log.error("failed to handle RabbitMQ message", ex);
        }
    }
}