type Price = {
  sku: string,
  unit: string,
  material: string,
  tservicio: string,
  price: number
};

type PriceList = {
  list: string,
  owner: string,
  currency: string,
  targets: string[],
  prices: Price[]
};

const API_URL = "http://localhost:8080";

export class PriceListsApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_URL) {
    this.baseUrl = baseUrl;
  }

  async getListsByOwnerAndTargets(owner: string, targets: string[]): Promise<{ lists: string[] }> {

    let targetsParam = '';
    for (const target of targets) {
      targetsParam += `&targets=${target}`;
    }

    const response = await fetch(`${this.baseUrl}/price-lists?owner=${owner}${targetsParam}`);
    if (!response.ok) {
      throw new Error(`Error fetching price list names: ${response.statusText}`);
    }
    return response.json();
  }

  async createPriceList(postData: PriceList): Promise<{results: string}> {

    const response = await fetch(`${this.baseUrl}/price-lists`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(postData),
    });
    if (!response.ok) {
      throw new Error(`Error creating price list: ${response.statusText}`);
    }
    return response.json();
  }
  
}
