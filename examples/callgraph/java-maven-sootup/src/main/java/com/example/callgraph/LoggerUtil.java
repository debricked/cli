package com.example.callgraph;

import org.apache.commons.lang3.StringUtils;

public final class LoggerUtil {
    private LoggerUtil() {
    }

    public static void log(String message) {
        String normalized = StringUtils.upperCase(message);
        System.out.println("[APP] " + normalized);
    }
}

