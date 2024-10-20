type Price = {
  sku: string;
  unit: string;
  material: string;
  tservicio: string;
  price: number;
};

type PriceListDTO = {
  label: string;
  created_at: number; // UNIX timestamp format
  currency: string;
};

const API_URL = "http://localhost:8080";

export class PriceListsApiClient {
  private baseUrl: string;
  private owner: string;

  constructor(owner: string, baseUrl: string = API_URL) {
    this.baseUrl = baseUrl;
    this.owner = owner;
  }

  // Private function to get lists by owner and targets
  private async getListsByOwnerAndTargets(targets: string[]): Promise<{ lists: string[] }> {

    const targetsParam = targets.map(target => `&targets=${target}`).join('');
    const response = await fetch(`${this.baseUrl}/price-lists?owner=${this.owner}${targetsParam}`);

    if (!response.ok) {
      throw new Error(`Error fetching price list names: ${response.statusText}`);
    }
    return response.json();
  }

  // Public function to fetch list (first found) for a single target
  public async fetchListForTarget(target: string): Promise<PriceListDTO> {

    const { lists } = await this.getListsByOwnerAndTargets([target]);

    const priceListDTO = lists.map(list => {
        const [currency, unixTime, label] = list.split('-');

        return {
          label,
          created_at: parseInt(unixTime, 10),
          currency,
        };
    });

    return priceListDTO[0]; // Return the first DTO
  }

  // Public function to fetch lists for multiple targets
  public async fetchListsForTargets(targets: string[]): Promise<PriceListDTO[]> {

    const { lists } = await this.getListsByOwnerAndTargets(targets);

    return lists.map(list => {
        const [currency, unixTime, label] = list.split('-');

        return {
          label,
          created_at: parseInt(unixTime, 10),
          currency,
        };
    });
  }

  public async createPriceListFromParams(
    list: string,
    currency: string,
    targets: string[],
    prices: Price[]
  ): Promise<{ results: string }> {

    const postData = {
      list: list,
      owner: this.owner,
      currency: currency,
      targets: targets,
      prices: prices,
    };

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

  public async addPrice(
    list: PriceListDTO,
    sku: string,
    unit: string,
    material: string,
    tservicio: string,
    price: number
  ): Promise<{ results: string }> {

    const postData = {
      list: `${list.currency}-${list.created_at}-${list.label}`, // Concatenate label, created_at, and currency
      sku,
      unit,
      material,
      tservicio,
      price,
    };

    const response = await fetch(`${this.baseUrl}/prices`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(postData),
    });

    if (!response.ok) {
      throw new Error(`Error adding price: ${response.statusText}`);
    }
    return response.json();
  }

  public async updatePrice(
    list: PriceListDTO,
    sku: string,
    unit: string,
    material: string,
    tservicio: string,
    price: number
  ): Promise<{ results: string }> {

    const postData = {
      list: `${list.label}-${list.created_at}-${list.currency}`, // Concatenate label, created_at, and currency
      sku,
      unit,
      material,
      tservicio,
      price,
    };

    const response = await fetch(`${this.baseUrl}/prices`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(postData),
    });

    if (!response.ok) {
      throw new Error(`Error updating price: ${response.statusText}`);
    }
    return response.json();
  }

  public async getPrice(
    list: PriceListDTO,
    sku: string,
    unit: string,
    material: string,
    tservicio: string
  ): Promise<Price> {

    const response = await fetch(
      `${this.baseUrl}/prices?list=${encodeURIComponent(`${list.label}-${list.created_at}-${list.currency}`)}&sku=${encodeURIComponent(sku)}&unit=${encodeURIComponent(unit)}&material=${encodeURIComponent(material)}&tservicio=${encodeURIComponent(tservicio)}`
    );

    if (!response.ok) {
      throw new Error(`Error fetching price: ${response.statusText}`);
    }
    return response.json();
  }
}
