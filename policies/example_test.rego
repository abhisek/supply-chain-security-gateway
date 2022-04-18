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
        "version": "c",
        "vulnerabilities": [],
        "licenses": []
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
        "version": "2.16.0",
        "vulnerabilities": [],
        "licenses": []
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
        "version": "2.17.2",
        "vulnerabilities": [],
        "licenses": []
      }
    }
  }

  result := violations with input as input
  count(result) = 0
}
