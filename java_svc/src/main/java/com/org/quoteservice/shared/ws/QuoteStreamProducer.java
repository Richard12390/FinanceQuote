package com.org.quoteservice.shared.ws;

import lombok.RequiredArgsConstructor;
import org.springframework.amqp.core.Message;
import org.springframework.beans.BeanUtils;
import org.springframework.rabbit.stream.producer.RabbitStreamTemplate;
import org.springframework.rabbit.stream.support.StreamMessageProperties;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class QuoteStreamProducer {

    private final RabbitStreamTemplate streamTemplate;

    public void publish(QuoteMessage quote) {
        streamTemplate.convertAndSend(quote, this::ensureStreamProperties);
    }

    private Message ensureStreamProperties(Message message) {
        if (message.getMessageProperties() instanceof StreamMessageProperties) {
            return message;
        }
        StreamMessageProperties streamProps = new StreamMessageProperties();
        BeanUtils.copyProperties(message.getMessageProperties(), streamProps);
        return new Message(message.getBody(), streamProps);
    }
}
