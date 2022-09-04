configurationPath="$HOME/.scs/pacman.env"

function loadConfigurationIfPresent() {
  if [ -f "$configurationPath" ]; then
    print_msg "Loading config from $configurationPath"
    source $configurationPath
  fi;
}

function interactiveConfiguration() {
  print_msg "Running interactive configuration"
  echo -n "Gateway Base URL (without trailing /): "
  read -r gatewayURL
  echo -n "Username: "
  read -r username
  echo -n "Password: "
  read -r password

  mkdir -p `dirname $configurationPath` 2>/dev/null
  cat > $configurationPath <<_EOF
GATEWAY_URL=$gatewayURL
GATEWAY_USERNAME=$username
GATEWAY_PASSWORD=$password
_EOF
}
