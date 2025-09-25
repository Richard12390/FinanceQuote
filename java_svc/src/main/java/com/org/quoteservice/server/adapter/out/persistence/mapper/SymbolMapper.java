package com.org.quoteservice.server.adapter.out.persistence.mapper;

import com.org.quoteservice.server.application.dto.SymbolUpsert;

public interface SymbolMapper {
    int upsertSymbol(SymbolUpsert cmd);

    Long lastInsertId();

    Integer findExchangeIdByCode(String code);

    Integer findCategoryIdByCode(String code);
}
