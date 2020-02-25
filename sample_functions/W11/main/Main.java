package main;

import lalb.*;

public class Main {
    public static lalb.Response lambda_handler() {
		double sum = 0;
        for (int i = 0; i <= 8000000; i++) {
			sum += Math.random();
		}

		return new lalb.Response(200, "Hello World" + Integer.toString((int)sum));
    }
}
