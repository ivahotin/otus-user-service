package com.example.billing

import org.springframework.dao.EmptyResultDataAccessException
import org.springframework.jdbc.core.namedparam.NamedParameterJdbcTemplate
import org.springframework.stereotype.Repository
import java.util.UUID

@Repository
class PaymentRepository(
        val jdbcTemplate: NamedParameterJdbcTemplate
) {

    fun withdraw(ownerId: String, amount: Long): Boolean {
        val rowsAffected = jdbcTemplate.update(
                "update billing_accounts set amount = amount - :amount where owner_id = :id::uuid and amount >= :amount",
                mapOf("id" to ownerId, "amount" to amount)
        )

        return rowsAffected > 0
    }

    fun replenish(ownerId: String, amount: Long) {
        jdbcTemplate.update(
                "update billing_accounts set amount = amount + :amount where owner_id = :id::uuid",
                mapOf("id" to ownerId, "amount" to amount)
        )
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