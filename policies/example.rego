package pdp

# Input format
# {"kind":"PolicyInput","version":{"major":1,"minor":0,"patch":0},
# "target":{"artefact":{"source":{"type":"maven2"},
# "group":"com.google.j2objc",
# "name":"j2objc-annotations",
# "version":"1.3",
# "vulnerabilities":[],"licenses":[]}}}

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

  msg := "Old and vulnerable version of log4j2 not allowed"
  code := 1002
}
