package user

const (
	name     = "Test Name"
	email    = "test@test.com"
	password = "Password@1234"

	emptyString = ""

	invalidPasswordOne   = "password@1234"
	invalidPasswordTwo   = "PASSWORD@1234"
	invalidPasswordThree = "Password@"
	invalidPasswordFour  = "Password1"
	invalidPasswordFive  = "Pa@1"

	passwordHash = "passwordHash"

	dirID  = "d87f0cb7-46c0-4501-83fa-e1c5ed5338e6"
	userID = "86d690dd-92a0-40ac-ad48-110c951e3cb8"

	saltLength = 86
	iterations = 4096
	keyLength  = 32

	pemString = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAiAAAABNlY2RzYS
1zaGEyLW5pc3RwMzg0AAAACG5pc3RwMzg0AAAAYQQqIb8uvVE0tZ7+XzpRaQhasWyJzgsW
+geO1C5tb+TEBSdtOqJD2z9Tss8p2EoEMqIyGeYwG0M4zD2KuL9qoF7RG+8V04U5FIBs0s
FgNbsG+gnRSvcdbwPbSaQuhICCi3QAAADw7ZHdZ+2R3WcAAAATZWNkc2Etc2hhMi1uaXN0
cDM4NAAAAAhuaXN0cDM4NAAAAGEEKiG/Lr1RNLWe/l86UWkIWrFsic4LFvoHjtQubW/kxA
UnbTqiQ9s/U7LPKdhKBDKiMhnmMBtDOMw9iri/aqBe0RvvFdOFORSAbNLBYDW7BvoJ0Ur3
HW8D20mkLoSAgot0AAAAMF+X6Raq5SFPv5Cd9+uLGzrqDcDSchOVJggL9emKfEvcb/M6XW
QI4c5bkbrw8Fm07AAAACZuaWtoaWxzb25pQE5pa2hpbHMtTWFjQm9vay1Qcm8tMi5sb2Nh
bAEC
-----END OPENSSH PRIVATE KEY-----`
)
