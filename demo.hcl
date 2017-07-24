name = "HCL"
params {
	url = "http://agilebits.com"
	tls {
		enabled = false
		pem = ["file://somewhere.pem", "file://more.pem"]
	}
}