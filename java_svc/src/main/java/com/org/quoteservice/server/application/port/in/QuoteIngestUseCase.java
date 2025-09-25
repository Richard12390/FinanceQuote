package com.org.quoteservice.server.application.port.in;

public interface QuoteIngestUseCase {
    void handleMessage(String subject, String jsonPayload);
}