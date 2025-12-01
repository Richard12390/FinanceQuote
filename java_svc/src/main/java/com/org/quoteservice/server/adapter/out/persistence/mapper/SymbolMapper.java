package com.org.quoteservice.server.adapter.out.persistence.mapper;

import com.org.quoteservice.server.application.dto.SymbolUpsert;

/*
import java.util.List;
import java.util.Map;
import com.org.quoteservice.server.application.dto.ListInstrumentsQuery;
*/

public interface SymbolMapper {
    int upsertSymbol(SymbolUpsert cmd);

    Long lastInsertId();

    Integer findExchangeIdByCode(String code);

    Integer findCategoryIdByCode(String code);

    /*
    List<Map<String, Object>> listInstruments(ListInstrumentsQuery query);
    */
}
