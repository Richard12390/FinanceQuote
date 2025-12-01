package com.org.quoteservice.server.domain.model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
@JsonIgnoreProperties(ignoreUnknown = true)
public class QuoteTick {
    private String symbol;
    private Double price;
    private Double bid;
    private Double ask;
    private Double bidSize;
    private Double askSize;
    private Double volume;
    private String marketState;
    private Long tsUnixMs;
    private String tsIso;
}
