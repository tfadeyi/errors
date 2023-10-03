#!/bin/sh

set -e

if ! command -v tar >/dev/null; then
	echo "Error: tar is required to install fyictl" 1>&2
	exit 1
fi

if ! command -v curl >/dev/null; then
    echo "Error: curl is required to install fyictl" 1>&2
    exit 1
fi


case $(uname -sm) in
"Darwin x86_64") target="fyictl-darwin-amd64" ;;
"Darwin arm64")  target="fyictl-darwin-arm64" ;;
"Linux x86_64")  target="fyictl-linux-amd64" ;;
"Linux aarch64") target="fyictl-linux-arm64" ;;
  *)
      echo "Error: Unsupported operating system or architecture: $(uname -sm)" 1>&2
      exit 1 ;;
esac
  target_file="errctl"


fyictl_uri="https://github.com/tfadeyi/errors/releases/latest/download/${target}.tar.gz"

fyictl_install="${FYICTL_INSTALL:-$HOME/.fyictl}"
bin_dir="$fyictl_install/bin"
bin="$bin_dir/$target_file"

if current_install=$( command -v fyictl ) && [ ! -x "$bin" ] && [ "$current_install" != "$bin" ]; then
    echo "failed to install fyictl to \"$bin\"" >&2
    echo "fyictl is already installed in another location: \"$current_install\"" >&2
    exit 1
fi

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$bin.tar.gz" "$fyictl_uri"
tar xfO "$bin.tar.gz" "$target/$target_file" > "$bin"
chmod +x "$bin"
rm "$bin.tar.gz"

echo "fyictl was installed successfully to $bin"
if command -v fyictl >/dev/null; then
	echo "Run 'fyictl --help' to get started"
else
	if [ "$SHELL" = "/bin/zsh" ] || [ "$ZSH_NAME" = "zsh" ]; then
        shell_profile=".zshrc"
    else
        shell_profile=".bashrc"
	fi
    echo
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export FYICTL_INSTALL=\"$fyictl_install\""
	echo "  export PATH=\"\$FYICTL_INSTALL/bin:\$PATH\""
    echo
	echo "And run \"source $HOME/$shell_profile\" to update your current shell"
    echo
fi
echo
echo "Checkout https://error.fyi/fyictl for more information"