package pdp

default allow = false

allow {
  count(violations) == 0
}

violations[{"message": msg, "code": code}] {
  input.kind != "PolicyInput"

  msg := "Input kind is unexpected in policy"
  code := 1000
}

violations[{"message": msg, "code": code}] {
  input.version.major != 1

  msg := "Input schema is not supported"
  code := 1001
}

violations[{"message": msg, "code": code}] {
  input.target.artefact.group = "org.apache.logging.log4j"
  input.target.artefact.name = "log4j"
  semver.compare(input.target.artefact.version, "2.17.0") = -1

  msg := "Old and vulnerable version of log4j2 is not allowed"
  code := 1002
}

violations[{"message": msg, "code": code}] {
  some i, j

  input.target.vulnerabilities[i].severity =
    data.UNACCEPTABLE_VULNERABILITIES[j]

  msg := sprintf("Vulnerabilities with %v severity blocked",
    [data.UNACCEPTABLE_VULNERABILITIES])
  code := 1003
}
