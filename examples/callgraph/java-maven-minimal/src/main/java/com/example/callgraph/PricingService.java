package com.example.callgraph;

public class PricingService {
    public double calculateTotal(int quantity, double unitPrice) {
        double subtotal = quantity * unitPrice;
        return applyDiscount(subtotal);
    }

    private double applyDiscount(double subtotal) {
        if (subtotal >= 50.0) {
            return subtotal * 0.9;
        }
        return subtotal;
    }
}

