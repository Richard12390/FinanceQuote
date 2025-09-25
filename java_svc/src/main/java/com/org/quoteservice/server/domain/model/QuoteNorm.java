package com.org.quoteservice.server.domain.model;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;
import java.util.Map;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@Data
@NoArgsConstructor
@AllArgsConstructor
@JsonIgnoreProperties(ignoreUnknown = true)
public class QuoteNorm {
    private String source;
    private String assetType;
    private String symbol;
    private String symbolOrigin;
    private String exchange;
    private String currency;
    private Double price;
    private Double open;
    private Double high;
    private Double low;
    private Double prevClose;
    private Double bid;
    private Double bidSize;
    private Double ask;
    private Double askSize;
    private Double volume;
    private String volumeKind;
    private String marketState;
    private Long tsUnixMs;
    private String tsLocal;
    private Map<String, Object> raw;
}
