package com.org.quoteservice.server.common.auth;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.security.Keys;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.Instant;
import java.util.Date;
import java.util.Map;
import javax.crypto.SecretKey;

public final class JwtUtil {
  private static final String SECRET = "asdf4568q9w5is2b6r8q9c5s8g45fpwa";
  private static final SecretKey SECRET_KEY = Keys.hmacShaKeyFor(SECRET.getBytes(StandardCharsets.UTF_8));
  private static final long EXPIRE_MILLIS = Duration.ofHours(24).toMillis();

  private JwtUtil() {
  }

  public static String generateToken(String subject) {
    return generateToken(subject, Map.of());
  }

  public static String generateToken(String subject, Map<String, Object> claims) {
    Date now = new Date();
    Date expiry = new Date(now.getTime() + EXPIRE_MILLIS);

    return Jwts.builder()
        .setSubject(subject)
        .addClaims(claims)
        .setIssuedAt(now)
        .setExpiration(expiry)
        .signWith(SignatureAlgorithm.HS256, SECRET_KEY)
        .compact();
  }

  public static Claims parseToken(String token) {
    return Jwts.parserBuilder()
        .setSigningKey(SECRET_KEY)
        .build()
        .parseClaimsJws(token)
        .getBody();
  }

  public static String parseSubject(String token) {
    return parseToken(token).getSubject();
  }

  public static String extractToken(String authHeader) {
    if (authHeader == null || !authHeader.startsWith("Bearer ")) {
      throw new IllegalArgumentException("Missing Bearer token");
    }
    return authHeader.substring(7);
  }

  public Instant getExpiration(String token) {
    Claims claims = parseToken(token);
    return claims.getExpiration().toInstant();
  }
}
