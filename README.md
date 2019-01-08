# CoffeeWithCloudySpartans
Developed an Online Coffee Store application using Mongo, Express, React, Node (MERN) and Go. Hosted the data service API modules developed in Go as microservices using Kubernetes.

## Application Architecture Diagram

![ArchitectureDiagram](/Architecture/Architecture.png)

### Microservices Architecture

### User Services

![UserServices](/Architecture/MongoDiagram.png)


### Catalog Services


![CatalogServices](/Architecture/SignUp_Login_architecture.png)


### Shopping Cart Services


![ShoppingCart](/Architecture/CartAPI_diagram.png)


### Payment Services


![PaymentServices](/Architecture/Payment_API.png)


Following are the core services in our application:
### login/sign-up service
A user will be able to sign-up and login using this service.

### catalog services (Administrator)
Items available will be displayed on the website.
Additionally maintenance of items within the catalog will be performed by users with admin privileges.

### shopping cart service
A user can add/remove multiple items to the cart and checkout.

### order processing service
Items present in the shopping cart are processed via a card payment or coupons.

## Team Members:
1. Mudambi Seshadri Srinivas 
2. Preethi Thimma Govarthanarajan
3. Abhishek Konduri
4. Hansraj Mathakar


