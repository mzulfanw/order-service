package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type ProductClient interface {
	GetProduct(id string) (*Product, error)
}

type productClient struct {
	baseURL string
	client  *http.Client
}

func NewProductClient(url string) ProductClient {
	return &productClient{
		baseURL: url,
		client:  &http.Client{},
	}
}

func (p *productClient) GetProduct(id string) (*Product, error) {
	url := fmt.Sprintf("%s/products/%s", p.baseURL, id)
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found (status %d)", resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product response: %w", err)
	}

	fmt.Printf("✅ Product fetched: %+v\n", product)
	return &product, nil
}
