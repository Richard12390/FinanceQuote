package com.org.quoteservice.shared.ws;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties("app.streams")
public record RabbitStreamProperties(
    String host,
    int port,
    String streamName,
    String consumerName,
    String username,
    String password
) {
}