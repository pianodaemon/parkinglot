# Sidecar Program for Price Management

## Overview

**Prices** is a side-car that runs alongside the main application to handle the management of prices pertaining to one or more lists. It is designed to work independently of the primary application logic while ensuring that all pricing-related operations are offloaded for better scalability, modularity, and performance.

### Key Design points:

-   **Price List**: Each list belongs to an owner.
-   **Target**: The list can be shared with multiple targets (customers or any other entity).
-   **Price**: Each price is associated with a unique hash and the corresponding tuple, ensuring that even if multiple items have the same attributes (like `sku` or `material`), they are uniquely identified.
   
## Key Features

- **RESTful API**: Exposes APIs to retrieve and update prices and price lists.
