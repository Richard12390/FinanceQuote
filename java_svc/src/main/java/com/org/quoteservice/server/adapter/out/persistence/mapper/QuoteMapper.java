package com.org.quoteservice.server.adapter.out.persistence.mapper;

import java.sql.Date;
import com.org.quoteservice.server.domain.model.QuoteNorm;

public interface QuoteMapper {
    int upsertDaily(long symbolId,
            QuoteNorm quote,
            String rawJson,
            Date tradeDate);
}
