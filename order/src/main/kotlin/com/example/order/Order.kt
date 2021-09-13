package com.example.order

import java.util.UUID

enum class OrderStatus {
    IN_PROGRESS, SUCCESS, FAILED
}

data class Order(
        val id: UUID,
        val ownerId: String,
        val price: Long,
        val status: OrderStatus,
        val version: Long
)
