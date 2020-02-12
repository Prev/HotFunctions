import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.ByteArrayOutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

public class Entry_HFRestServer {

    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress(8080), 0);
        server.createContext("/", new MyHandler());
        server.setExecutor(null); // creates a default executor
        server.start();

        System.out.println("Service running at 8080");
    }

    static class MyHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange t) throws IOException {
            Runtime runtime = Runtime.getRuntime();
            Process process;

            try {
                process = runtime.exec("java -cp .:lib/gson-2.8.6.jar Entry_HF");
                process.waitFor();

            } catch(Exception e) {
                String response = "Error";
                t.sendResponseHeaders(200, response.length());

                OutputStream os = t.getResponseBody();
                os.write(response.getBytes());
                os.close();
                return;

            }
            InputStream is = process.getInputStream();
            ByteArrayOutputStream buffer = new ByteArrayOutputStream();
            int nRead;
            byte[] data = new byte[16384];
            while ((nRead = is.read(data, 0, data.length)) != -1) {
              buffer.write(data, 0, nRead);
            }

            byte[] bytes = buffer.toByteArray();

            t.sendResponseHeaders(200, bytes.length);

            OutputStream os = t.getResponseBody();
            os.write(bytes);
            os.close();
        }
    }
}
