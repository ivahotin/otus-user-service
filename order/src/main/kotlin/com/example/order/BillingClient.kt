package com.example.order

import org.springframework.beans.factory.annotation.Value
import org.springframework.http.HttpEntity
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpMethod
import org.springframework.stereotype.Component
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.RestTemplate
import java.io.Serializable

data class Payment(val amount: Long): Serializable

@Component
class BillingClient(
    @Value("\${billing.url}") private val billingUrl: String,
    @Value("\${billing_port}") private val billingPort: String
) {

    private val restTemplate = RestTemplate()

    fun payForOrder(userId: String, amount: Long): Boolean {
        val httpHeaders = HttpHeaders()
        httpHeaders.add("x-user-id", userId)
        val requestBody = HttpEntity<Payment>(Payment(amount), httpHeaders)

        return try {
            restTemplate.exchange(
                    "http://$billingUrl:$billingPort/payments/withdrawals",
                    HttpMethod.PUT,
                    requestBody,
                    Any::class.java
            )
            true
        } catch (exc: HttpClientErrorException) {
            return false
        }
    }
}