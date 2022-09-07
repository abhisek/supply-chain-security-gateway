package pdp

test_violation_kind_fail {
  input := {
    "kind": "Invalid",
    "version": { "major": 1, "minor": 0, "patch": 0 },
    "target": {
      "artefact": {
        "source": { "type": "maven2" },
        "group": "a",
        "name": "b",
        "version": "c"
      }
    }
  }

  result := violations with input as input
  count(result) = 1
}

test_log4j_old_version_fail {
  input := {
    "kind": "PolicyInput",
    "version": { "major": 1, "minor": 0, "patch": 0 },
    "target": {
      "artefact": {
        "source": { "type": "maven2" },
        "group": "org.apache.logging.log4j",
        "name": "log4j",
        "version": "2.16.0"
      }
    }
  }

  result := violations with input as input
  count(result) = 1
}

test_log4j_old_version_pass {
  input := {
    "kind": "PolicyInput",
    "version": { "major": 1, "minor": 0, "patch": 0 },
    "target": {
      "artefact": {
        "source": { "type": "maven2" },
        "group": "org.apache.logging.log4j",
        "name": "log4j",
        "version": "2.17.2"
      }
    }
  }

  result := violations with input as input
  count(result) = 0
}

test_fail_on_critical_vuln {
  input := {
    "kind": "PolicyInput",
    "version": { "major": 1, "minor": 0, "patch": 0 },
    "target": {
      "artefact": {
        "source": { "type": "maven2" },
        "group": "org.example",
        "name": "random",
        "version": "1.33.7",
      },
      "vulnerabilities": [
        { "severity": "CRITICAL" }
      ]
    }
  }

  result := violations with input as input
  count(result) = 1
}

test_pass_on_low_vuln {
  input := {
    "kind": "PolicyInput",
    "version": { "major": 1, "minor": 0, "patch": 0 },
    "target": {
      "artefact": {
        "source": { "type": "maven2" },
        "group": "org.example",
        "name": "random",
        "version": "1.33.7",
      },
      "vulnerabilities": [
        { "severity": "LOW" }
      ]
    }
  }

  result := violations with input as input
  count(result) = 0
}

test_fail_private_namespace {
  input := {
    "kind": "PolicyInput",
    "version": { "major": 1, "minor": 0, "patch": 0 },
    "target": {
      "artefact": {
        "source": { "type": "maven2" },
        "group": "org.example.private.lib",
        "name": "random",
        "version": "1.33.7",
      },
      "vulnerabilities": [
        { "severity": "LOW" }
      ]
    }
  }

  result := violations with input as input
  count(result) = 1
}
