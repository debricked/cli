package com.example.callgraph;

public class App {
    public static void main(String[] args) {
        PricingService pricingService = new PricingService();
        OrderService orderService = new OrderService(pricingService);

        String confirmation = orderService.placeOrder("book", 2, 24.99);
        LoggerUtil.log(confirmation);
    }
}

