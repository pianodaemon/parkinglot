import { PriceListsApiClient } from "./apiClient.ts";

const owner = "maestro_limpio";
const client = new PriceListsApiClient(owner);

async function runClient() {
  try {
    // -------------------------------------------
    // Create a price list
    console.log("\nCreating a new price list...");

    const list = "inv-2060";
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
        unit: "litros",
        material: "vidrio",
        tservicio: "limpia",
        price: 1500.00
      },
    ];

    const results = await client.createPriceListFromParams(list, currency, targets, prices);
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

  } catch (error) {
    console.error(error);
  }
}

// Learn more at https://docs.deno.com/runtime/manual/examples/module_metadata#concepts
if (import.meta.main) {
  runClient();
}
