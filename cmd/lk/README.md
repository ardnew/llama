# 🥾 lk

<p align="center">
  <br>
  <img src=".github/images/demo.gif" width="600" alt="lk demo">
  <br>
</p>

**lk** — a terminal navigator.

Taken from [antonmedv/walk](https://github.com/antonmedv/walk):

> Why another terminal navigator? I wanted something simple and minimalistic. Something to help me with faster navigation in the filesystem; a `cd` and `ls`replacement. So I build \[`lk`\]. It allows for quick navigation with fuzzy searching, `cd` integration is quite simple. And you can open `vim` right from \[`lk`\]. That's it.

## Install

```
go install github.com/ardnew/walk/v2/cmd/lk@latest
```

Or download [prebuild binaries](https://github.com/ardnew/walk/releases).

Put the next function into the `.bashrc` or a similar config:

<table>
<tr>
  <th> Bash/Zsh </th>
  <th> Fish </th>
  <th> PowerShell </th>
</tr>
<tr>
<td>

```bash
function clk {
  cd "$(lk "$@")"
}
```

</td>
<td>

```fish
function clk
  set loc (lk $argv); and cd $loc;
end
```

</td>
<td>

```powershell
function clk() {
  cd $(lk $args)
}
```

</td>
</tr>
</table>


Now use `clk` command to start walking.

## Usage

| Key binding      | Description        |
|------------------|--------------------|
| `Arrows`, `hjkl` | Move cursor        |
| `Enter`          | Enter directory    |
| `Backspace`      | Exit directory     |
| `Space`          | Toggle preview     |
| `Esc`, `q`       | Exit with cd       |
| `Ctrl+c`         | Exit without cd    |
| `/`              | Fuzzy search       |
| `dd`             | Delete file or dir |
| `y`              | yank current dir   |

The `EDITOR` or `LK_EDITOR` environment variable used for opening files from lk.

```bash
export EDITOR=vim
```

### Preview mode

Press `Space` to toggle preview mode.

<img src=".github/images/preview-mode.gif" width="600" alt="Walk Preview Mode">

### Delete file or directory

Press `dd` to delete file or directory. Press `u` to undo.

<img src=".github/images/rm-demo.gif" width="600" alt="Walk Deletes a File">

### Display icons

Install [Nerd Fonts](https://www.nerdfonts.com) and add `--icons` flag.

<img src=".github/images/demo-icons.gif" width="600" alt="Walk Icons Support">

### Image preview

No additional setup is required.

<img src=".github/images/images-mode.gif" width="600" alt="Walk Image Preview">

## Become a sponsor

Every line of code in my repositories 📖 signifies my unwavering commitment to open source 💡. Your support 🤝 ensures these projects keep thriving, innovating, and benefiting all 💼. If my work has ever resonated 🎵 or helped you, kindly consider showing love ❤️ by sponsoring. [**🚀 Sponsor Me Today! 🚀**](https://github.com/sponsors/antonmedv)

## License

[MIT](LICENSE)
