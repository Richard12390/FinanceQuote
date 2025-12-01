package com.org.quoteservice.server.common.auth;

import java.io.IOException;
import java.util.List;

import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import io.jsonwebtoken.Claims;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@Component
public class JwtAuthenticationFilter extends OncePerRequestFilter {

  @Override
  protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response,
      FilterChain filterChain) throws ServletException, IOException {

    String authHeader = request.getHeader("Authorization");
    if (authHeader != null && authHeader.startsWith("Bearer ")) {
      try {
        String token = authHeader.substring(7);
        Claims claims = JwtUtil.parseToken(token);
        Long userId = claims.get("userId", Long.class);
        String account = claims.getSubject();

        JwtUserPrincipal principal = new JwtUserPrincipal(userId, account);
        UsernamePasswordAuthenticationToken authentication = new UsernamePasswordAuthenticationToken(principal, null,
            List.of());
        SecurityContextHolder.getContext().setAuthentication(authentication);
      } catch (Exception ex) {
        response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Invalid token");
        return;
      }
    }
    filterChain.doFilter(request, response);
  }

  @Override
  protected boolean shouldNotFilter(HttpServletRequest request) {
    String path = request.getRequestURI();
    return path.startsWith("/auth") || HttpMethod.OPTIONS.matches(request.getMethod());
  }
}
