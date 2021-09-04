package com.example.billing

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*

@RestController
class BillingController(private val repository: PaymentRepository) {

    @PutMapping("/payments/withdrawals")
    fun withdraw(@RequestHeader("x-user-id") userId: String, @RequestBody payment: PaymentRequest): ResponseEntity<*> {
        val isPaid = repository.withdraw(userId, payment.amount)
        return if (isPaid) ResponseEntity.ok().build<Any>() else ResponseEntity.status(HttpStatus.CONFLICT).build<Any>()
    }

    @PutMapping("/payments/replenishments")
    fun replenish(@RequestHeader("x-user-id") userId: String, @RequestBody payment: PaymentRequest): ResponseEntity<*> {
        repository.replenish(userId, payment.amount)

        return ResponseEntity.ok().build<Any?>()
    }

    @GetMapping("/payments")
    fun getUserBilling(@RequestHeader("x-user-id") userId: String): BillingAccount? {
        return repository.getBillingAccountByOwnerId(userId)
    }
}