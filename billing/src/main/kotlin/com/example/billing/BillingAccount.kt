package com.example.billing

import java.util.UUID

data class BillingAccount(
        val ownerId: UUID,
        val amount: Long
)
