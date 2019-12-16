package main;

public class Main {
    public static String lambda_handler() {
		double sum = 0;
        for (int i = 0; i <= 6000000; i++) {
			sum += Math.random();
		}

		return "Hello World" + Integer.toString((int)sum);
    }
}
