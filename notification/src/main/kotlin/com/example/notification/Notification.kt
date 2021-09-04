package com.example.notification

import java.util.UUID

data class Notification(
        val id: Int,
        val ownerId: UUID,
        val orderId: UUID,
        val price: Int,
        val isSuccess: Boolean
)