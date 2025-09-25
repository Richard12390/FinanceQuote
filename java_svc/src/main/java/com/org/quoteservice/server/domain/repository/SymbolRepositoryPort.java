package com.org.quoteservice.server.domain.repository;

import java.util.Optional;
import com.org.quoteservice.server.domain.model.Symbol;

public interface SymbolRepositoryPort {
    long upsertSymbol(Symbol symbol);

    Optional<Integer> findExchangeIdByCode(String code);

    Optional<Integer> findCategoryIdByCode(String code);
}