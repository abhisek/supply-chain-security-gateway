mavenScriptSrc="$BASEDIR/data/maven-settings.xml"
mavenScriptDst="$HOME/.m2/settings.xml"

function setupMaven() {
  gatewayURL=$1
  username=$2
  password=$3

  gatewayMavenCentralURL="$gatewayURL$MAVEN_CENTRAL_ROUTE"
  mkdir -p `dirname $mavenScriptDst` 2>/dev/null

  if [ -f "$mavenScriptDst" ]; then
    echo "[WARN] Overwriting $mavenScriptDst"
  fi;

  cat $mavenScriptSrc | \
    sed "s,{{GATEWAY_USERNAME}},$username," | \
    sed "s,{{GATEWAY_PASSWORD}},$password," | \
    sed "s,{{GATEWAY_MAVEN_CENTRAL_URL}},$gatewayMavenCentralURL," \
    > $mavenScriptDst
}

setupMaven $GATEWAY_URL $GATEWAY_USERNAME $GATEWAY_PASSWORD
