sequenceDiagram

User->>+OrderService: POST /orders {userId}
OrderService->>OrderService: Create order
OrderService->>+BillingService: PUT /payments {orderId}
BillingService -->>- OrderService: 201 Created
OrderService ->>+ MessageBroker: publish
OrderService -->>- User: 201 Created
MessageBroker -->>+ NotificationService: consume
Note right of MessageBroker: OrderCreated
NotificationService ->>+ UserService: GET /users/{userId}
UserService -->>- NotificationService: 200
NotificationService ->>- NotificationService: sending notification
