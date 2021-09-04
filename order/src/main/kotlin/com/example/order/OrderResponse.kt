package com.example.order

import java.io.Serializable
import java.util.UUID

data class OrderResponse(
        val id: UUID,
        val is_success: Boolean
): Serializable