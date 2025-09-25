package com.org.quoteservice.server.domain.repository;

import java.sql.Date;
import com.org.quoteservice.server.domain.model.QuoteNorm;

public interface QuoteRepositoryPort {
    void upsertDaily(long symbolId, QuoteNorm quote, String rawJson, Date tradeDate);
}