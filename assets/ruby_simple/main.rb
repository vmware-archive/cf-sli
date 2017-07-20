require "webrick"
STDOUT.sync = true
STDERR.sync = true

Thread.new do
  while true do
    STDOUT.puts "Tick: #{Time.now.to_i}"
    sleep 1
  end
end

class EnvServlet < WEBrick::HTTPServlet::AbstractServlet
    def do_GET (request, response)
      if request.unparsed_uri == '/'
        response.status = 200
        response.body = <<-RESPONSE
Healthy
It just needed to be restarted!
My application metadata: #{ENV['VCAP_APPLICATION']}
My port: #{ENV['PORT']}
My custom env variable: #{ENV['CUSTOM_VAR']}
RESPONSE
      else
        response.status = 404
      end
    end
end

class LogServlet < WEBrick::HTTPServlet::AbstractServlet
    def do_GET (request, response)
      m = request.unparsed_uri.match(/^\/log\/(.*)/)
      if m && m[1]
        message = m[1]
        response.status = 200
        STDOUT.puts(message)
        response.body = "logged #{message} to STDOUT"
      else
        response.status = 404
      end
    end
end

server = WEBrick::HTTPServer.new(:Port => ENV.fetch('PORT', '8080'))

server.mount "/log", LogServlet
server.mount "/", EnvServlet

trap("INT") {
    server.shutdown
}

server.start
