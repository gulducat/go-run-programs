program "consul" {
  command = "consul agent -dev"
  check   = "consul members"
}
  
program "nomad" {
	command = "nomad agent -dev"
	check   = "nomad node status"
	# seconds = 1 # to test check failures
}
  
program "vault" {
  command = "vault server -dev"
  check   = "vault status"
  env = {
    # default client is https, but vault server -dev is http
    VAULT_ADDR = "http://localhost:8200"
  }
}
