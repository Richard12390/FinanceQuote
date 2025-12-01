package com.org.quoteservice.server.common.auth;

public record JwtUserPrincipal(Long userId, String account) {
}
