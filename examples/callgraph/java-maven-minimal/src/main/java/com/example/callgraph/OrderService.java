package com.example.callgraph;

public class OrderService {
    private final PricingService pricingService;

    public OrderService(PricingService pricingService) {
        this.pricingService = pricingService;
    }

    public String placeOrder(String item, int quantity, double unitPrice) {
        double total = pricingService.calculateTotal(quantity, unitPrice);
        return "Order placed: " + item + " x" + quantity + " total=" + total;
    }
}

