#!/bin/sh

set -e

if ! command -v tar >/dev/null; then
	echo "Error: tar is required to install errctl" 1>&2
	exit 1
fi

if ! command -v curl >/dev/null; then
    echo "Error: curl is required to install errctl" 1>&2
    exit 1
fi


case $(uname -sm) in
"Darwin x86_64") target="errctl-darwin-amd64" ;;
"Darwin arm64")  target="errctl-darwin-arm64" ;;
  "Linux x86_64")  target="errctl-linux-amd64" ;;
"Linux aarch64") target="errctl-linux-arm64" ;;
  *)
      echo "Error: Unsupported operating system or architecture: $(uname -sm)" 1>&2
      exit 1 ;;
esac
  target_file="errctl"


errctl_uri="https://github.com/tfadeyi/errors/releases/latest/download/${target}.tar.gz"

errctl_install="${ERRCTL_INSTALL:-$HOME/.sloscribe}"
bin_dir="$errctl_install/bin"
bin="$bin_dir/$target_file"

if current_install=$( command -v sloscribe ) && [ ! -x "$bin" ] && [ "$current_install" != "$bin" ]; then
    echo "failed to install errctl to \"$bin\"" >&2
    echo "errctl is already installed in another location: \"$current_install\"" >&2
    exit 1
fi

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$bin.tar.gz" "$errctl_uri"
tar xfO "$bin.tar.gz" "$target/$target_file" > "$bin"
chmod +x "$bin"
rm "$bin.tar.gz"

echo "errctl was installed successfully to $bin"
if command -v errctl >/dev/null; then
	echo "Run 'errctl --help' to get started"
else
	if [ "$SHELL" = "/bin/zsh" ] || [ "$ZSH_NAME" = "zsh" ]; then
        shell_profile=".zshrc"
    else
        shell_profile=".bashrc"
	fi
    echo
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export ERRCTL_INSTALL=\"$errctl_install\""
	echo "  export PATH=\"\$ERRCTL_INSTALL/bin:\$PATH\""
    echo
	echo "And run \"source $HOME/$shell_profile\" to update your current shell"
    echo
fi
echo
echo "Checkout https://error.fyi/errctl for more information"