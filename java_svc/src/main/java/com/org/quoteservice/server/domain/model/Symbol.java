package com.org.quoteservice.server.domain.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Symbol {
    private String source;
    private int categoryId;
    private String symbol;
    private String symbolOrigin;
    private Integer exchangeId;
    private String currency;
    private int isActive;
}