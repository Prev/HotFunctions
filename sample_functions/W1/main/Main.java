package main;

class Result {
	public int statusCode;
	public String body;

	Result(int statusCode, String body) {
		this.statusCode = statusCode;
		this.body = body;
	}
}

public class Main {
    public static Result lambda_handler() {
		double sum = 0;
        for (int i = 0; i <= 6000000; i++) {
			sum += Math.random();
		}

		return new Result(200, "Hello World" + Integer.toString((int)sum));
    }
}
