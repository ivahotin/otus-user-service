package com.example.billing

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*

@RestController
class BillingController(private val repository: PaymentRepository) {

    @PutMapping("/payments/withdrawals")
    fun withdraw(
            @RequestHeader("x-user-id") userId: String,
            @RequestHeader("idempotency-key") idempotencyKey: String,
            @RequestBody payment: PaymentRequest
    ): ResponseEntity<*> {
        val operationResult = try {
            repository.withdraw(idempotencyKey, userId, payment.amount)
        } catch (exc: Throwable) {
            return when (exc) {
                is InsufficientAmount -> ResponseEntity.status(HttpStatus.PRECONDITION_FAILED).build<Any?>()
                else -> ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build<Any?>()
            }
        }

        return when (operationResult) {
            is PaymentMade -> ResponseEntity.ok().build<Any>()
            is PaymentWasMadeBefore -> ResponseEntity.status(HttpStatus.CONFLICT).build<Any?>()
            is InsufficientAmount -> ResponseEntity.status(HttpStatus.PRECONDITION_FAILED).build<Any?>()
        }
    }

    @PutMapping("/payments/replenishments")
    fun replenish(
            @RequestHeader("x-user-id") userId: String,
            @RequestHeader("idempotency-key") idempotencyKey: String,
            @RequestBody payment: PaymentRequest
    ): ResponseEntity<*> {
        return when (repository.replenish(idempotencyKey, userId, payment.amount)) {
            is PaymentMade -> ResponseEntity.ok().build<Any?>()
            is PaymentWasMadeBefore -> ResponseEntity.status(HttpStatus.CONFLICT).build<Any?>()
            else -> ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).build<Any?>()
        }
    }

    @GetMapping("/payments")
    fun getUserBilling(@RequestHeader("x-user-id") userId: String): BillingAccount? {
        return repository.getBillingAccountByOwnerId(userId)
    }
}