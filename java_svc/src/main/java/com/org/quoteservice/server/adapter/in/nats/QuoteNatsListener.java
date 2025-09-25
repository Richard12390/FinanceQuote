package com.org.quoteservice.server.adapter.in.nats;

import io.nats.client.Connection;
import io.nats.client.Dispatcher;
import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import com.org.quoteservice.server.application.port.in.QuoteIngestUseCase;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@RequiredArgsConstructor
public class QuoteNatsListener {
    private final Connection connection;
    private final QuoteIngestUseCase ingestUseCase;

    @PostConstruct
    public void subscribe() {
        Dispatcher dispatcher = connection.createDispatcher(message -> {
            try {
                ingestUseCase.handleMessage(message.getSubject(), new String(message.getData()));
            } catch (Exception e) {
                log.error("failed to handle NATS message", e);
            }
        });
        dispatcher.subscribe("quotes.>");
        log.info("NATS dispatcher subscribed to quotes.>");
    }
}