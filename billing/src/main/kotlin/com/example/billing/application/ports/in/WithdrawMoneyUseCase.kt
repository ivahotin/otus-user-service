package com.example.billing.application.ports.`in`

interface WithdrawMoneyUseCase {
    fun withdrawMoney(userId: String, amount: Long): Boolean
}