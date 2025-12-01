package com.org.quoteservice.shared.ws;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties("app.stomp")
public record StompBrokerProperties(
    String host,
    int port,
    String virtualHost,
    String systemLogin,
    String systemPasscode,
    String clientLogin,
    String clientPasscode
) {
}