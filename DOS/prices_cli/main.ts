import { PriceListsApiClient } from "./apiClient.ts";

const owner = "maestro_limpio";
const client = new PriceListsApiClient(owner);
const listNameStr = (list: { label: string, created_at: number, currency: string }) =>
  `${list.currency}-${list.created_at}-${list.label}`;

async function runClient() {
  try {
    // -------------------------------------------
    // Create a price list
    console.log("\nCreating a new price list...");

    const listLabel = "invierno 2060";
    const currency = "MXN";
    const targets = [
      "coca",
      "pepsi"
    ];
    const prices = [
      {
        sku: "123123-222",
        unit: "m3",
        material: "madera",
        tservicio: "recoleccion",
        price: 99.99
      },
      {
        sku: "2345-987999",
        unit: "kg",
        material: "vidrio",
        tservicio: "limpia",
        price: 1500.00
      },
    ];

    const results = await client.createPriceListFromParams(listLabel, currency, targets, prices);
    console.log(results);

    // -------------------------------------------
    // Get a price list name by owner and target
    console.log(`\nFetching a price list name by owner ${owner} and target ${targets[0]}...`);
    const listName = await client.fetchListForTarget(targets[0]);
    console.log(listName);

    // -------------------------------------------
    // Get price list names by owner and targets
    console.log(`\nFetching all price list names by owner ${owner} and targets ${targets}...`);
    const listNames = await client.fetchListsForTargets(targets);
    console.log(listNames);

    // -------------------------------------------
    // Add a price to a specific price list (the fetched one in previous lines: listName)
    console.log(`\nAdding a price to a specific price list ${listNameStr(listName)}...`);

    const sku = "123-4512-22";
    const unit = "kg";
    const material = "cascajo";
    const tservicio = "recoleccion";
    const price = 450.10;

    const addPriceRes = await client.addPrice(listName, sku, unit, material, tservicio, price);
    console.log(addPriceRes);

    // -------------------------------------------
    // Get a price from a specific price list
    console.log(`\nGetting a price from a specific price list ${listNameStr(listName)}...`);
    let priceRes = await client.getPrice(listName, sku, unit, material, tservicio);
    console.log(priceRes);

    // -------------------------------------------
    // Update a price from a specific price list
    console.log(`\nUpdating a price from a specific price list ${listNameStr(listName)}...`);
    const updatePriceRes = await client.updatePrice(listName, sku, unit, material, tservicio, 9999.99);
    console.log(updatePriceRes);

    // -------------------------------------------
    // Get a recently updated price
    console.log(`\nGetting a recently updated price from a specific price list ${listNameStr(listName)}...`);
    priceRes = await client.getPrice(listName, sku, unit, material, tservicio);
    console.log(priceRes);

  } catch (error) {
    console.error(error);
  }
}

// Learn more at https://docs.deno.com/runtime/manual/examples/module_metadata#concepts
if (import.meta.main) {
  runClient();
}
