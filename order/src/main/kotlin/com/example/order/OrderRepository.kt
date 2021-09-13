package com.example.order

import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import java.lang.RuntimeException
import java.util.UUID

typealias Version = Long
typealias CreationStatus = Boolean

@Repository
class OrderRepository(private val jdbcTemplate: JdbcTemplate) {

    fun createOrder(order: Order): Pair<Version, CreationStatus> {
        val isThereConflict = jdbcTemplate.update(
                "insert into orders (id, owner_id, price, status, version) values (?::uuid, ?::uuid, ?, ?, ?) on conflict (owner_id, version) do nothing",
                order.id,
                order.ownerId,
                order.price,
                order.status.toString(),
                order.version
        ) == 0

        return if (isThereConflict) {
            val version = jdbcTemplate.queryForObject(
                    "select version as latest_version from orders where owner_id = ?::uuid order by version desc limit 1",
                    arrayOf(order.ownerId)
            ) {
                rs, _ -> rs.getLong("latest_version")
            } ?: throw RuntimeException("Something goes wrong")
            version to false
        } else {
            order.version to true
        }
    }

    fun updateOrderStatusById(orderId: UUID, status: OrderStatus): Int {
        return jdbcTemplate.update(
                "update orders set status = ? where id = ?::uuid",
                status.toString(),
                orderId
        )
    }

    fun getOrdersByOwnerId(ownerId: String): List<Order> {
        return jdbcTemplate.query(
                "select id, owner_id, price, status, version from orders where owner_id = ?::uuid order by version desc",
                arrayOf(ownerId)
        ) {
            rs, _ -> Order(
                id = UUID.fromString(rs.getString("id")),
                ownerId = rs.getString("owner_id"),
                price = rs.getLong("price"),
                status = OrderStatus.valueOf(rs.getString("status")),
                version = rs.getLong("version")
            )
        }
    }
}