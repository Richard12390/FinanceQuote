package com.org.quoteservice.server.domain.service;

import java.util.Locale;
import lombok.RequiredArgsConstructor;
import com.org.quoteservice.server.domain.repository.SymbolRepositoryPort;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class SymbolDomainService {
    private static final int DEFAULT_CATEGORY_ID = 1;

    private final SymbolRepositoryPort symbolRepository;

    @Cacheable(cacheNames = "symbolCategory", key = "#assetType != null ? #assetType.toLowerCase() : 'default'")
    public int resolveCategoryId(String assetType) {
        if (assetType == null) {
            return DEFAULT_CATEGORY_ID;
        }
        return symbolRepository.findCategoryIdByCode(assetType.toLowerCase(Locale.ROOT))
                .orElse(DEFAULT_CATEGORY_ID);
    }
}
