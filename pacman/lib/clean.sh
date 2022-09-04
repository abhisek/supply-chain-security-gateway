function remove_file() {
  file=$1

  echo "[WARN] Removing file: $file"
  rm -f $file
}

remove_file "$HOME/$GRADLE_INIT_SCRIPT_PATH"
remove_file "$HOME/.m2/settings.xml"
