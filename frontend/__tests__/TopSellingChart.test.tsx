import React from "react";
import { render, screen } from "@testing-library/react";
import TopSellingChart from "@/components/TopSellingChart";
import type { ProductSaleInfo } from "@/lib/types";

const items: ProductSaleInfo[] = [
  { id: 1, name: "Widget", sku: "SKU-1", quantity: 10, revenue: 100 },
  { id: 2, name: "Gadget", sku: "SKU-2", quantity: 5, revenue: 50 },
];

it("renders a row per product with quantity and revenue", () => {
  render(<TopSellingChart items={items} />);

  expect(screen.getByText("Widget")).toBeInTheDocument();
  expect(screen.getByText("Gadget")).toBeInTheDocument();
  expect(screen.getByText("10 sold")).toBeInTheDocument();
  expect(screen.getByText("5 sold")).toBeInTheDocument();
  expect(screen.getByText("$100.00 revenue")).toBeInTheDocument();
});

it("scales the best seller to 100% and shows an accessible label", () => {
  render(<TopSellingChart items={items} />);

  const top = screen.getByLabelText("Widget: 10 sold");
  expect(top.firstChild).toHaveStyle({ width: "100%" });
});

it("shows an empty state with no items", () => {
  render(<TopSellingChart items={[]} />);
  expect(screen.getByText("No sales yet.")).toBeInTheDocument();
});