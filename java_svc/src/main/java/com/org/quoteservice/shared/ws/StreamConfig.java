package com.org.quoteservice.shared.ws;

import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@EnableConfigurationProperties({RabbitStreamProperties.class, StompBrokerProperties.class})
public class StreamConfig {
}