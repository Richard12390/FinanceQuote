package com.org.quoteservice.server.adapter.out.persistence.repository;

import java.util.Optional;
import lombok.RequiredArgsConstructor;
import com.org.quoteservice.server.adapter.out.persistence.mapper.SymbolMapper;
import com.org.quoteservice.server.application.dto.SymbolUpsert;
import com.org.quoteservice.server.application.port.out.SymbolRepositoryPort;
import com.org.quoteservice.server.domain.model.Symbol;
import org.springframework.stereotype.Repository;

@Repository
@RequiredArgsConstructor
public class MyBatisSymbolRepository implements SymbolRepositoryPort {
    private final SymbolMapper symbolMapper;

    @Override
    public long upsertSymbol(Symbol symbol) {
        SymbolUpsert command = new SymbolUpsert(
                symbol.getSource(),
                symbol.getCategoryId(),
                symbol.getSymbol(),
                symbol.getSymbolOrigin(),
                symbol.getExchangeId(),
                symbol.getCurrency(),
                symbol.getIsActive());
        symbolMapper.upsertSymbol(command);
        Long sid = symbolMapper.lastInsertId();
        return sid == null ? 0L : sid;
    }

    @Override
    public Optional<Integer> findExchangeIdByCode(String code) {
        if (code == null || code.isEmpty()) {
            return Optional.empty();
        }
        return Optional.ofNullable(symbolMapper.findExchangeIdByCode(code));
    }

    @Override
    public Optional<Integer> findCategoryIdByCode(String code) {
        if (code == null || code.isEmpty()) {
            return Optional.empty();
        }
        return Optional.ofNullable(symbolMapper.findCategoryIdByCode(code));
    }
}
