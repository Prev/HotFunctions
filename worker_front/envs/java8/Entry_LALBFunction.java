import com.google.gson.Gson;
import com.google.gson.JsonObject;
import main.Main;

public class Entry_LALBFunction {
    public static void main(String args[]) {

		long startTime = System.currentTimeMillis();
		String out = Main.lambda_handler();
		long endTime = System.currentTimeMillis();

		Gson gson = new Gson();
		JsonObject object = new JsonObject();
		object.addProperty("startTime", startTime);
		object.addProperty("endTime", endTime);
		object.addProperty("body", out);
		String jsonString = gson.toJson(object);

		System.out.printf("-----%d-----=%s\n", jsonString.length(), jsonString);
    }
}
