package com.example.billing

import org.springframework.dao.DuplicateKeyException
import org.springframework.dao.EmptyResultDataAccessException
import org.springframework.jdbc.core.namedparam.NamedParameterJdbcTemplate
import org.springframework.stereotype.Repository
import org.springframework.transaction.annotation.Transactional
import java.util.UUID

@Repository
class PaymentRepository(
        val jdbcTemplate: NamedParameterJdbcTemplate
) {

    @Transactional(rollbackFor = [InsufficientAmount::class])
    fun withdraw(idempotencyKey: String, ownerId: String, amount: Long): PaymentOperationResult {
        try {
            jdbcTemplate.update(
                    "insert into transactions (idempotency_key, created_at, amount, is_cancelled) values (:key::uuid, now(), -:amount, false)",
                    mapOf("key" to idempotencyKey, "amount" to amount)
            )
        } catch (exc: DuplicateKeyException) {
            return PaymentWasMadeBefore
        }

        val rowsAffected = jdbcTemplate.update(
                "update billing_accounts set amount = amount - :amount where owner_id = :id::uuid and amount >= :amount",
                mapOf("amount" to amount, "id" to ownerId)
        )

        if (rowsAffected > 0) {
            return PaymentMade
        }

        throw InsufficientAmount
    }

    @Transactional
    fun replenish(idempotencyKey: String, ownerId: String, amount: Long): PaymentOperationResult {
        try {
            jdbcTemplate.update(
                    "insert into transactions (idempotency_key, created_at, amount, is_cancelled) values (:key::uuid, now(), -:amount, false)",
                    mapOf("key" to idempotencyKey, "amount" to amount)
            )
        } catch (exc: DuplicateKeyException) {
            return PaymentWasMadeBefore
        }

        jdbcTemplate.update(
                "update billing_accounts set amount = amount + :amount where owner_id = :id::uuid",
                    mapOf("id" to ownerId, "amount" to amount, "idempotencyKey" to idempotencyKey)
            )
        return PaymentMade
    }

    fun getBillingAccountByOwnerId(ownerId: String): BillingAccount? {
        return try {
            jdbcTemplate.queryForObject(
                    "select * from billing_accounts where owner_id = :id::uuid",
                    mapOf("id" to ownerId)
            ) {
                rc, _ ->
                BillingAccount(UUID.fromString(rc.getString("owner_id")), rc.getLong("amount"))
            }
        } catch (exc: EmptyResultDataAccessException) {
            throw BillingNotFoundException()
        }
    }
}