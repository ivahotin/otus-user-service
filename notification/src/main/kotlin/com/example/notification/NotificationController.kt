package com.example.notification

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestHeader
import org.springframework.web.bind.annotation.RestController

@RestController
class NotificationController(private val notificationRepository: NotificationRepository) {

    @GetMapping("/notifications")
    fun getNotificationByUserId(@RequestHeader("x-user-id") userId: String): List<NotificationMessage> {
        return notificationRepository
                .getNotificationByOwnerId(userId)
                .map {
                    when (it.isSuccess) {
                        true -> NotificationMessage(
                                "Dear ${it.ownerId}. Order ${it.orderId} has been successfully processed")
                        false -> NotificationMessage(
                                "Dear ${it.ownerId}. Order ${it.orderId} couldn't been processed")
                    }
                }
    }
}