gradleScriptSrc="$BASEDIR/data/plugin.gradle"
gradleScriptDst="$HOME/$GRADLE_INIT_SCRIPT_PATH"

function setupGradle() {
  gatewayURL=$1
  username=$2
  password=$3

  gatewayPluginPortalURL="$gatewayURL$GRADLE_PLUGINS_ROUTE"
  gatewayMavenCentralURL="$gatewayURL$MAVEN_CENTRAL_ROUTE"

  echo sed "s/{{GATEWAY_USERNAME}}/$username/"

  mkdir -p `dirname $gradleScriptDst` 2>/dev/null
  cat $gradleScriptSrc | \
    sed "s,{{GATEWAY_USERNAME}},$username," | \
    sed "s,{{GATEWAY_PASSWORD}},$password," | \
    sed "s,{{GATEWAY_GRADLE_PLUGIN_URL}},$gatewayPluginPortalURL," | \
    sed "s,{{GATEWAY_MAVEN_CENTRAL_URL}},$gatewayMavenCentralURL," \
    > $gradleScriptDst
}

setupGradle $GATEWAY_URL $GATEWAY_USERNAME $GATEWAY_PASSWORD
