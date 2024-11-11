# <img src="https://i.ibb.co/RNpWkb6/flxvwr.png" alt="logo" width="40" style="vertical-align: middle;"/> flxvwr

[![Go Report Card](https://goreportcard.com/badge/github.com/justjcurtis/flxvwr)](https://goreportcard.com/report/github.com/justjcurtis/flxvwr)
[![License](https://img.shields.io/github/license/justjcurtis/flxvwr)](https://github.com/justjcurtis/flxvwr/blob/main/Licence.md)

This project is a simple, cross-platform image viewer written in Go using the [Fyne](https://fyne.io/) toolkit. flxvwr provides a minimal user interface, keyboard shortcuts for navigation, and smooth performance across platforms including Windows, macOS, and Linux.

## Features

- **Cross-Platform Support**: Runs on Windows, macOS, and Linux.
- **Simple UI**: Minimalistic interface for a clean viewing experience.
- **Keyboard Shortcuts**: Navigate through images and control the viewer with intuitive key bindings.
- **Responsive**: Built with Fyne, enabling a responsive and native look on all platforms.
- **Performance**: Smooth performance for viewing images with minimal overhead with no image count limit.
- **Customisable**: Adjust brightness, contrast, zoom, pan and rotation.
- **Playlist**: Drag&drop images, directories or playlist files to import for viewing.
    - **Recursive Scan**: Recursively scan directories and playlist files for images, directories and playlist files to import.
    - **Clear Playlist**: Clear the current playlist to start fresh.
- **Shuffle Mode**: Toggle shuffle mode to view images in a random order.
- **Delay**: Adjust the delay between images for a custom viewing experience.

## Screenshots

| [![Welcome Screen](https://i.ibb.co/H78vt6f/welcome.png)](https://ibb.co/jhXKymP) | [![Feature Screen](https://i.ibb.co/0s1hXYV/image.png)](https://ibb.co/jJqg53Z) | [![Settings Screen](https://i.ibb.co/JkkgQ4G/settings.png)](https://ibb.co/vzzRYGT) |
|---------------------------------------------|---------------------------------------------|---------------------------------------------|
| Welcome Screen                              | Image Screen                                | Settings Screen                             |

## Keyboard Shortcuts

- **Q**: Quit flxvwr
- **Esc**: Quit all instances of flxvwr
- **C**: Clear the current playlist
- **Space**: Play/Pause the current playlist
- **Left/Right Arrow**: Go to the previous/next item
- **Up/Down Arrow**: Adjust delay by +/- 1 second
- **S**: Toggle shuffle mode
- **+/-**: Zoom in/out
- **H/L/J/K**: Pan left/right/up/down
- **[ / ]**: Rotate image
- **R**: Reset image
- **B/N**: Adjust brightness up/down
- **V/M**: Adjust contrast up/down
- **Shift+...**: Smaller increments for adjustments and movements
- **Shift+1 to 9**: Add current image to a numbered playlist slot
- **1 to 9**: Switch to a numbered playlist slot
- **X**: Remove current image from the playlist
- **E**: Export the current playlist to a `.txt` file
- **F1 | ?**: Show/hide shortcuts
- **/**: Open settings

## Installation

### Via Package Manager

#### Arch Linux [(AUR Package)](https://aur.archlinux.org/packages/flxvwr-bin)

1. **Install flxvwr**:
    ```bash
    yay -S flxvwr-bin
    ```
2. **Run flxvwr**:
    ```bash
    flxvwr
    ```
    **or**

    `Search for and start flxvwr from your application launcher.`

#### MacOS and Windows

macOS (Homebrew Cask) and Windows (Chocolatey Package) coming soon.

### Via Release Download (Windows, macOS, Linux)

1. **Download the latest release**:

    Download the latest release for your os from the [releases page](https://github.com/justjcurtis/flxvwr/releases/latest).

2. Install the release:
    - **Windows**: Run the installer and follow the instructions.
    - **macOS**: Open the `.dmg` file and drag the app to your Applications folder.
    - **Linux**: Install the tar.xz file as you would any other package. (eg. extract, cd into the directory, sudo make install)

### Via `go install` (Linux, macOS, Windows)

1. **Get Go**:

    Go [here](https://go.dev/doc/install) to install Go.

2. **Ensure Go paths are set**:
    Add the following to your `.bashrc` or `.zshrc`:
    ```bash
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    ```

3. **Install flxvwr**:
    ```bash
    go install github.com/justjcurtis/flxvwr@latest
    ```

4. **Run flxvwr**
   ```bash
   flxvwr
   ```


## Usage

Launch flxvwr, drag&drop images, directories and/or playlist files into the window, and start viewing. Use the keyboard shortcuts to navigate between images and manage the viewer. Directories and playlist files are recursively scanned for images, and you can clear the current playlist using `C` or toggle shuffle mode using `S` for a random viewing experience.

## Playlist Files

- Playlist files are `.txt` files containing a newline separated list of filepaths. 
- These filepaths can be images, directories and other playlist files. flxvwr will recursively scan directories and playlist files for images to import.
- Any lines in the playlist file that are not valid filepaths that point to an existing file/directory will be ignored.
- Even if a playlist file refers to itself or it's parent directory, for example, flxvwr will handle this gracefully and prevent infinite recursion.

## Supported File Types

- **JPEG/JPG**
- **PNG**
- **TXT**
- **Directory**

## Roadmap

- **Publish Releases**: Publish releases to ~~AUR~~, Homebrew, and Chocolatey.
- **Customisable Key Bindings**: Allow users to customise key bindings for navigation and controls.
- **Image Metadata**: Display image metadata such as resolution, file size, and format.
- **Image Sorting**: Sort images by name, date, or size.
- **Image Filtering**: Filter images by file type or resolution.
- **Plugins**: Add support for Lua plugins to extend flxvwr's functionality.
    - **Package Manager**: Implement a package manager for installing and managing plugins.
- **Thumbnail View**: Display a grid of image thumbnails for quicker navigation.
- **UI Tests**: Write tests to ensure flxvwr functions as expected.

## Contributing

Contributions are welcome! Feel free to submit a pull request to suggest improvements or additional features.

## License

This project is licensed under the MIT License.
