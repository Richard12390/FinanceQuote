package com.org.quoteservice.server.adapter.out.persistence.repository;

import java.sql.Date;
import lombok.RequiredArgsConstructor;
import com.org.quoteservice.server.adapter.out.persistence.mapper.QuoteMapper;
import com.org.quoteservice.server.application.port.out.QuoteRepositoryPort;
import com.org.quoteservice.server.domain.model.QuoteNorm;
import org.springframework.stereotype.Repository;

@Repository
@RequiredArgsConstructor
public class MyBatisQuoteRepository implements QuoteRepositoryPort {
    private final QuoteMapper quoteMapper;

    @Override
    public void upsertDaily(long symbolId, QuoteNorm quote, String rawJson, Date tradeDate) {
        quoteMapper.upsertDaily(symbolId, quote, rawJson, tradeDate);
    }
}