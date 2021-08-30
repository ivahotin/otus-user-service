sequenceDiagram

participant User
participant Gateway
participant MessageBroker
participant OrderService
participant BillingService
participant NotificationService
participant UserService

User ->>+ Gateway: POST /orders{userId, price}
Gateway ->>+ MessageBroker: publish
Note right of MessageBroker: OrderCreationRequested
MessageBroker -->>+ OrderService: consume
OrderService ->> MessageBroker: publish
Note right of MessageBroker: OrderPaymentRequested
MessageBroker -->>+ BillingService: consume
Note left of BillingService: OrderPaymentRequested
BillingService ->>+ OrderService: GET /orders/{orderId}
OrderService -->>- BillingService: 200
BillingService ->>- MessageBroker: publish
Note right of MessageBroker: OrderPaid
MessageBroker -->> OrderService: consume
Note left of OrderService: OrderPaid
OrderService ->>- MessageBroker: publish
Note right of MessageBroker: OrderCreated
MessageBroker ->>- Gateway: consume
Gateway -->> User: Response
MessageBroker -->>+ NotificationService: consume
Note left of NotificationService: OrderCreated
NotificationService ->>+ UserService: GET /users/{userId}
UserService -->>- NotificationService: 200
NotificationService ->>+ OrderService: GET /orders/{orderId}
OrderService -->>- NotificationService: 200
NotificationService ->> NotificationService: sending email
NotificationService ->>- MessageBroker: publish
Note right of MessageBroker: NotificationSent
