package com.example.billing.application.ports.`in`

interface ReplenishmentMoneyUseCase {
    fun replenishMoney(userId: String, amount: Long)
}