package com.example.order

import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestHeader
import org.springframework.web.bind.annotation.RestController
import java.util.UUID

@RestController
class OrderController(private val orderRepository: OrderRepository, private val billingClient: BillingClient) {

    @PostMapping("/orders")
    fun createOrder(@RequestHeader("x-user-id") userId: String, @RequestBody order: OrderRequest): OrderResponse {
        val orderId = UUID.randomUUID()
        val isPaid = billingClient.payForOrder(userId, order.price)
        orderRepository.createOrder(Order(orderId, userId, order.price, isPaid))
        return OrderResponse(orderId, isPaid)
    }
}