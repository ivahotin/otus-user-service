package com.example.order

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*
import java.util.UUID

@RestController
class OrderController(private val orderRepository: OrderRepository, private val billingClient: BillingClient) {

    @PostMapping("/orders")
    fun createOrder(
            @RequestHeader("x-user-id") userId: String,
            @RequestHeader("If-Match") version: Long,
            @RequestBody orderRequest: OrderRequest
    ): ResponseEntity<*> {
        val orderId = UUID.randomUUID()
        val order = Order(orderId, userId, orderRequest.price, status = OrderStatus.IN_PROGRESS, version = version)
        val (latestVersion, successfullyCreated) = orderRepository.createOrder(order)
        if (!successfullyCreated) {
            return ResponseEntity
                    .status(HttpStatus.CONFLICT)
                    .header("ETag", latestVersion.toString())
                    .build<Any?>()
        }
        val isPaid = billingClient.payForOrder(userId, orderId, order.price)
        val orderStatus = if (isPaid) OrderStatus.SUCCESS else OrderStatus.FAILED
        orderRepository.updateOrderStatusById(orderId, orderStatus)
        return ResponseEntity
                .status(HttpStatus.OK)
                .header("ETag", latestVersion.toString())
                .body(Order(orderId, userId, orderRequest.price, orderStatus, version))
    }

    @GetMapping("/orders")
    fun getOrders(@RequestHeader("x-user-id") userId: String): ResponseEntity<List<Order>> {
        val orders = orderRepository.getOrdersByOwnerId(ownerId = userId)
        return ResponseEntity
                .status(HttpStatus.OK)
                .header("ETag", orders.firstOrNull()?.version.toString() ?: "0")
                .body(orders)
    }
}