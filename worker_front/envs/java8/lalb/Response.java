package lalb;

public class Response {
	public int statusCode;
	public String body;

	public Response(int statusCode, String body) {
		this.statusCode = statusCode;
		this.body = body;
	}
}
