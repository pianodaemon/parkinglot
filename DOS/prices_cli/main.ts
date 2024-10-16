import { PriceListsApiClient } from "./apiClient.ts";

const client = new PriceListsApiClient();

async function runClient() {
  try {
    // -------------------------------------------
    // Get price list's names by owner and targets
    const owner = 'chilli-willie'
    const targets = ['coca', 'pepsi']

    console.log(`\nFetching all price list names by owner ${owner} and targets ${targets}...`);
    const listNames = await client.getListsByOwnerAndTargets(owner, targets);
    console.log(listNames);

    // -------------------------------------------
    // Create a price list
    console.log("\nCreating a new price list...");
    const priceList = {
      list: "inv-2060",
      owner: "siempre-limpio",
      currency: "MXN",
      targets: [
        "coca",
        "pepsi"
      ],
      prices: [
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
      ]
    };
    const results = await client.createPriceList(priceList);
    console.log(results);

  } catch (error) {
    console.error(error);
  }
}

// Learn more at https://docs.deno.com/runtime/manual/examples/module_metadata#concepts
if (import.meta.main) {
  runClient();
}
