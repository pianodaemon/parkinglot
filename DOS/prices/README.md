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
    "list": "winter-1728533139",  # <--- this is a unique index
    "owner": "viajes Ponchito"
}
```

```js
# Collection "targets"

{
    "list": "winter-1728533139",
    "target": "pepsi"
}

{
    "list": "winter-2024-1728533139",
    "target": "coca"
}
```

```js
# Collection "prices"

{
      "tuple": {
            "list": "winter-1728533139",
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


## Key Features

- **RESTful API**: Exposes APIs to retrieve and update prices and price lists.
