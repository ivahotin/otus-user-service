package com.example.notification

import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import java.util.UUID

@Repository
class NotificationRepository(private val jdbcTemplate: JdbcTemplate) {

    fun getNotificationByOwnerId(ownerId: String): List<Notification> {
        return jdbcTemplate.query(
                "select id, order_id, owner_id, price, is_success from notifications where owner_id = ?::uuid",
                arrayOf(ownerId)
        ) {
            rc, _ -> Notification(
                id = rc.getInt("id"),
                orderId = UUID.fromString(rc.getString("order_id")),
                ownerId = UUID.fromString(rc.getString("owner_id")),
                price = rc.getInt("price"),
                isSuccess = rc.getBoolean("is_success")
            )
        }
    }
}