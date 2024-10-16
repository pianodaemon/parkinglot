# Sidecar Program for Price Management

## Overview

**Prices** is a side-car that runs alongside the main application to handle the management of prices pertaining to one or more lists. It is designed to work independently of the primary application logic while ensuring that all pricing-related operations are offloaded for better scalability, modularity, and performance.

### Key Design points:

-   **The Lists**: Each list belongs to an owner.
-   **The targets**: The list can be shared with multiple targets (customers or any other entity).
-   **The prices**: Each price is associated with a unique hash and the corresponding tuple, ensuring that even if multiple items have the same attributes (like `sku` or `material`), they are uniquely identified.

### How the data is stored and addressed

The next json fragments elaborate on one of the several lists owned by "viajes Ponchito" ( This is basically how data is being structured within MongoDB ). 

```js
# Collection "price_lists"

{
    "list": "EUR-1728533139-winter",  # <--- this is a unique index
    "owner": "viajes Ponchito"
}
```

```js
# Collection "targets"

{
    "list": "EUR-1728533139-winter",
    "target": "pepsi"
}

{
    "list": "EUR-1728533139-winter",
    "target": "coca"
}
```

```js
# Collection "prices"

{
      "tuple": {
            "list": "EUR-1728533139-winter",
            "sku": "1254-545-66",
            "unit": "m3",
            "material": "madera",
            "tservicio": "limpia"
      },
     "hash": "6f5902ac237024bdd0c176cb93063dc4",  # <--- this is a unique index 
     "price": 15.50
}
```

Based on the structure of the MongoDB collections above, they represent a relationship between price lists, their owners, and the targets (e.g., customers) for which the price lists are shared. Each price is uniquely identified using an MD5 hash derived from a tuple of values such as `sku`, `unit`, `material`, etc.

Here's a breakdown of the collections:

1.  **List Collection**:
    
    -   Each price list is owned by an entity, like "viajes Ponchito."
    -   The `list` field is a unique identifier for that specific price list.
2.  **Target Collection**:
    
    -   This links the price list to one or more targets, such as "pepsi" or "coca."
    -   Each target is associated with the same list.
3.  **Price Collection**:
    
    -   This stores the actual prices for items, represented by a `tuple` that includes fields like `sku`, `unit`, `material`, and `tservicio`.
    -   A unique hash is generated for each price, ensuring that the combination of tuple elements is unique.
    -   The `price` field stores the actual price of the item.


## How to use it

The following examples demonstrate how to create, update, retrieve, and manage price lists and their associated prices via the sidecar program's REST API.

### Create price list

This command creates a new price list for a specific owner with associated targets and prices.

```sh

pianodaemon@LAPTOP-4RSVIK4C:~$ curl --location 'localhost:8080/price-lists' \
--header 'Content-Type: application/json' \
--data '{
    "list": "winter",
    "owner": "viajes Ponchito",
    "currency": "EUR",
    "targets": [
        "coca",
        "pepsi"
    ],
    "prices": [
        {
            "sku": "123123-222",
            "unit": "m3",
            "material": "madera",
            "tservicio": "recoleccion",
            "price": 99.99
        },
        {
            "sku": "2345-987999",
            "unit": "litros",
            "material": "vidrio",
            "tservicio": "limpia",
            "price": 1500.00
        },
        {
            "sku": "84738-382777",
            "unit": "m3",
            "material": "aluminio",
            "tservicio": "recoleccion",
            "price": 2.99
        },
        {
            "sku": "84738-3820000",
            "unit": "litros",
            "material": "energon",
            "tservicio": "quema",
            "price": 111.11
        }
    ]
}'
```
**Explanation**:

-   This `POST` request creates a new price list called `winter`, owned by `viajes Ponchito` for `EUR` currency.
-   The list is shared with two targets: `coca` and `pepsi`.
-   Four prices are added, each corresponding to a different combination of attributes like `sku`, `unit`, `material`, and `tservicio`.
-   Each price is assigned a specific value under the `price` field, allowing for granular control over the pricing of different items.

### Add newer price to the list

This command adds a new price to an existing price list.

```sh

pianodaemon@LAPTOP-4RSVIK4C:~$ curl --location 'localhost:8080/prices' \
--data '{
    "list": "EUR-1728533139-winter",
    "sku": "84738-382888",
    "unit": "m3",
    "material": "radiactivo",
    "tservicio": "destruccion",
    "price": 2222.99
}'
```
**Explanation**:

-   This `POST` request adds a new price for the item with `sku` 84738-382888 to the price list `winter-1728533139`.
-   The item has attributes such as `unit` (m3), `material` (radiactivo), and `tservicio` (destruccion).
-   The price for this item is set to `2222.99`.

### Update price from a list

This command updates an existing price within a price list.

```sh

pianodaemon@LAPTOP-4RSVIK4C:~$ curl --location --request PUT 'localhost:8080/prices' \
--header 'Content-Type: application/json' \
--data '{
    "list": "EUR-1728533139-winter",
    "sku": "84738-382777",
    "unit": "m3",
    "material": "aluminio",
    "tservicio": "recoleccion",
    "price": 10099.99
}'
```
**Explanation**:

-   This `PUT` request updates the price for the item with `sku` 84738-382777 in the list `winter-1728533139`.
-   The price for this combination of `unit`, `material`, and `tservicio` is updated to `10099.99`.

### Retrieve price by tuple

This command retrieves the price for a specific item based on its tuple.

```sh

pianodaemon@LAPTOP-4RSVIK4C:~$ curl --location  \
'localhost:8080/prices?list=EUR-1728533139-winter&sku=84738-382777&unit=m3&material=aluminio&tservicio=recoleccion'
```

**Explanation**:

-   This `GET` request retrieves the price for the item in the price list `winter-1728533139` with the following attributes:
    -   `sku`: 84738-382777
    -   `unit`: m3
    -   `material`: aluminio
    -   `tservicio`: recoleccion
-   The price for this specific combination of attributes will be returned.

### Get lists from intersection of owner and targets

This command retrieves the price lists that belong to a specific owner and are shared with specific targets.

```sh

pianodaemon@LAPTOP-4RSVIK4C:~$ curl --location  \
'localhost:8080/price-lists?owner=viajes%20Ponchito&targets=coca&targets=pepsi'
```

**Explanation**:

-   This `GET` request retrieves all price lists owned by `viajes Ponchito` that are shared with both `coca` and `pepsi`.
-   The `owner` and `targets` parameters are used to filter the lists.
