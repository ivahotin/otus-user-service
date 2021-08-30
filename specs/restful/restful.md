sequenceDiagram

User->>+OrderService: POST /orders {userId}
OrderService->>OrderService: Create order
OrderService->>+BillingService: PUT /payments {orderId}
BillingService -->>- OrderService: 201 Created
OrderService->>+NotificationService: POST /notifications {template_id, type, context}
NotificationService -->>- OrderService: 201 Created
OrderService -->>- User: 201 Created