package com.org.quoteservice.server.adapter.out.persistence.mapper;

import java.util.List;

import org.apache.ibatis.annotations.Param;

import com.org.quoteservice.server.domain.model.Favorite;

public interface FavoriteMapper {

  List<Favorite> getByUserId(@Param("userId") long userId);

  void upsert(Favorite favorite);

  void delete(@Param("userId") long userId, @Param("symbol") String symbol);
}

