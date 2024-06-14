# Orders-Microservice:
This repo was created to apply my knowledge of microservices and GoLang.

# Project description:
The project demonstrates microservices concepts using GoLang + PostgreSQL + gRPC + DataDog + Protocol buffers. 



# Workflows:
- Placing order:
  - Will place order 
  - Will publish OrderCreated, consumed by restaurant service to process an order.
  - Will publish OrderStatusChanged, consumed by notification service to notify customers about order state changes.

- Update order status:
  - Update order status
  - Publish OrderStatusChanged   


