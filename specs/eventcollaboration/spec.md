sequenceDiagram

participant User
participant Gateway
participant MessageBroker
participant OrderService
participant BillingService
participant NotificationService
participant UserService

User ->>+ Gateway: POST /orders {userId, price}
Gateway ->>+ MessageBroker: publish
Note right of MessageBroker: OrderCreationRequested
MessageBroker -->>+ OrderService: consume
OrderService ->> MessageBroker: publish
Note right of MessageBroker: OrderPaymentRequested
MessageBroker -->>+ BillingService: consume
Note left of BillingService: OrderPaymentRequested
BillingService ->>- MessageBroker: publish
Note right of MessageBroker: OrderPaid
MessageBroker -->> OrderService: consume
Note left of OrderService: OrderPaid
OrderService ->>- MessageBroker: publish
Note right of MessageBroker: OrderCreated
MessageBroker ->>- Gateway: consume
Gateway -->> User: Response
UserService -->>+ MessageBroker: publish
Note left of NotificationService: UserUpdated
MessageBroker ->>- NotificationService: consume
Note left of NotificationService: UserUpdated
MessageBroker -->>+ NotificationService: consume
Note left of NotificationService: OrderCreated
NotificationService ->>- NotificationService: sending email