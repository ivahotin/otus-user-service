package com.example.order

import java.util.UUID

data class Order(
        val id: UUID,
        val ownerId: String,
        val price: Long,
        val isSuccess: Boolean
)
