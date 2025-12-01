package com.org.quoteservice.shared.ws;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Configuration;
import org.springframework.messaging.simp.config.MessageBrokerRegistry;
import org.springframework.web.socket.config.annotation.EnableWebSocketMessageBroker;
import org.springframework.web.socket.config.annotation.StompEndpointRegistry;
import org.springframework.web.socket.config.annotation.WebSocketMessageBrokerConfigurer;

@Configuration
@EnableWebSocketMessageBroker
@RequiredArgsConstructor
public class WebSocketConfig implements WebSocketMessageBrokerConfigurer {

    private final StompBrokerProperties brokerProperties;

    @Override
    public void configureMessageBroker(MessageBrokerRegistry registry) {
        registry.enableStompBrokerRelay("/topic")
                .setRelayHost(brokerProperties.host())
                .setRelayPort(brokerProperties.port())
                .setVirtualHost(brokerProperties.virtualHost())
                .setSystemLogin(brokerProperties.systemLogin())
                .setSystemPasscode(brokerProperties.systemPasscode())
                .setClientLogin(brokerProperties.clientLogin())
                .setClientPasscode(brokerProperties.clientPasscode());
        registry.setApplicationDestinationPrefixes("/app");
    }

    @Override
    public void registerStompEndpoints(StompEndpointRegistry registry) {
        registry.addEndpoint("/ws")
                .setAllowedOriginPatterns("*")
                .withSockJS();
    }
}