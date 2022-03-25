program "web1" {
  command = "test-server start 8081"
  check   = "test-server check 8081"
}
program "web2" {
  command = "test-server start 8082"
  check   = "test-server check 8082"
}
program "web3" {
  command = "test-server start 8083"
  check   = "test-server check 8083"
  seconds = 5
}
