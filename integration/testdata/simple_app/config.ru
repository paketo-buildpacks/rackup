class HelloWorld
  def call(env)
    [200, {"Content-Type" => "text/plain", "Content-Length" => "12"}, ["Hello world!"]]
  end
end

run HelloWorld.new
