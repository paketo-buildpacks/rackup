# This file is used by Rack-based servers to start the application.
#\ -p 3000

class HelloWorld
  def call(env)
    [200, {"Content-Type" => "text/plain", "Content-Length" => "12"}, ["Hello world!"]]
  end
end

run HelloWorld.new
