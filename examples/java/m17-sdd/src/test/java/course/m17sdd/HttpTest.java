package course.m17sdd;

import org.junit.jupiter.api.Test;
import java.net.*;
import java.net.http.*;
import java.time.Clock;
import java.util.Map;
import static org.junit.jupiter.api.Assertions.*;

class HttpTest {
    @Test void httpContract() throws Exception {
        var server = Server.start(new Access(true, Clock.systemUTC()), 0);
        try {
            String base = "http://127.0.0.1:" + server.getAddress().getPort();
            call(base,"POST","/requests","","{\"resource\":\"handbook\",\"hours\":1}",401);
            for (String body : new String[]{"{}","[]","{\"resource\":\"other\",\"hours\":1}","{\"resource\":\"handbook\",\"hours\":1.5}","{\"resource\":\"handbook\",\"hours\":\"1\"}","{\"resource\":\"handbook\",\"hours\":null}","{\"resource\":\"handbook\",\"hours\":false}","{\"resource\":\"handbook\"}","{\"resource\":\"handbook\",\"hours\":1,\"role\":\"approver\"}","{\"resource\":\"handbook\",\"hours\":1} {}"})
                call(base,"POST","/requests","alice",body,400);
            String id=(String)call(base,"POST","/requests","alice","{\"resource\":\"handbook\",\"hours\":1}",201).get("id");
            call(base,"GET","/requests/"+id,"bob","",403);
            call(base,"GET","/requests/missing","alice","",404);
            call(base,"POST","/requests/"+id+"/decision","bob","{\"decision\":\"approved\"}",403);
            call(base,"POST","/requests/"+id+"/decision","carol","{\"decision\":\"maybe\"}",400);
            call(base,"POST","/requests/"+id+"/decision","carol","{\"decision\":\"approved\"}",200);
            assertEquals(true,call(base,"GET","/access?user=alice&resource=handbook","alice","",200).get("allowed"));
            call(base,"GET","/access?user=alice&resource=handbook","bob","",403);
            call(base,"POST","/requests/"+id+"/revoke","carol","{\"reason\":\"done\"}",200);
            assertEquals(false,call(base,"GET","/access?user=alice&resource=handbook","alice","",200).get("allowed"));
            call(base,"GET","/requests/"+id+"/extra","alice","",404);
        } finally { server.stop(0); }
    }
    @Test void baselineHttp() throws Exception {
        var server=Server.start(new Access(false,Clock.systemUTC()),0);
        String base="http://127.0.0.1:"+server.getAddress().getPort();
        try {
            for(String body:new String[]{"{\"resource\":\"handbook\"}","{\"resource\":\"runbook\",\"hours\":0}"}) {
                String id=(String)call(base,"POST","/requests","alice",body,201).get("id");
                call(base,"POST","/requests/"+id+"/revoke","carol","{\"reason\":\"done\"}",404);
            }
            call(base,"POST","/requests","alice","{\"resource\":\"handbook\",\"hours\":1}",400);
        } finally {server.stop(0);}
    }
    Map<String,Object> call(String base,String method,String path,String user,String body,int status) throws Exception {
        var request=HttpRequest.newBuilder(URI.create(base+path)).header("X-User",user)
            .method(method,HttpRequest.BodyPublishers.ofString(body)).build();
        var response=HttpClient.newHttpClient().send(request,HttpResponse.BodyHandlers.ofString());
        assertEquals(status,response.statusCode(),response.body()); return Json.parse(response.body());
    }
}
