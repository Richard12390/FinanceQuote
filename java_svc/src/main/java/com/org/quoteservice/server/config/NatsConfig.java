package com.org.quoteservice.server.config;

import io.nats.client.Connection;
import io.nats.client.Nats;
import io.nats.client.Options;
import io.nats.client.impl.ErrorListenerLoggerImpl;
import java.time.Duration;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class NatsConfig {
    @Bean
    public Connection nats(@Value("") String url) throws Exception {
        Options opts = new Options.Builder()
                .server(url)
                .connectionTimeout(Duration.ofSeconds(5))
                .maxReconnects(-1)
                .reconnectWait(Duration.ofSeconds(2))
                .errorListener(new ErrorListenerLoggerImpl())
                .build();
        return Nats.connect(opts);
    }
}