package com.org.quoteservice.server.application.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class SymbolUpsert {
    private String source;
    private int categoryId;
    private String symbol;
    private String symbolOrigin;
    private Integer exchangeId;
    private String currency;
    private int isActive;
}