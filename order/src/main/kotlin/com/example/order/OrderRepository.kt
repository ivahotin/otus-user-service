package com.example.order

import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository

@Repository
class OrderRepository(private val jdbcTemplate: JdbcTemplate) {

    fun createOrder(order: Order) {
        jdbcTemplate.update(
                "insert into orders (id, owner_id, price, is_success) values (?::uuid, ?::uuid, ?, ?)",
                order.id,
                order.ownerId,
                order.price,
                order.isSuccess
        )
    }
}