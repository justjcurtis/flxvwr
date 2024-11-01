# flxvwr

This project is a simple, cross-platform image viewer written in Go using the [Fyne](https://fyne.io/) toolkit. The application provides a minimal user interface, keyboard shortcuts for navigation, and smooth performance across platforms including Windows, macOS, and Linux.

## Features

- **Cross-Platform Support**: Runs on Windows, macOS, and Linux.
- **Simple UI**: Minimalistic interface for a clean viewing experience.
- **Keyboard Shortcuts**: Navigate through images and control the viewer with intuitive key bindings.
- **Responsive**: Built with Fyne, enabling a responsive and native look on all platforms.

## Screenshots

| [![Welcome Screen](https://i.ibb.co/H78vt6f/welcome.png)](https://ibb.co/jhXKymP) | [![Feature Screen](https://i.ibb.co/0s1hXYV/image.png)](https://ibb.co/jJqg53Z) | [![Settings Screen](https://i.ibb.co/JkkgQ4G/settings.png)](https://ibb.co/vzzRYGT) |
|---------------------------------------------|---------------------------------------------|---------------------------------------------|
| Welcome Screen                              | Image Screen                                | Settings Screen                             |

## Keyboard Shortcuts

- **Q**: Quit the application
- **Esc**: Quit all instances of the application
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
- **/**: Open settings

## Installation

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

4. **Run the application**
   ```bash
   flxvwr
   ```

## Usage

Launch the application, drag&drop images or directories into the window, and start viewing. Use the keyboard shortcuts to navigate between images and manage the viewer. Directories are recursively scanned for images, and you can clear the current playlist using `C` or toggle shuffle mode using `S` for a random viewing experience.

## Roadmap

- **Release Scripts**: Add windows to release script & package macOS releases as a `.dmg` with a `.app` inside.
- **Publish Releases**: Publish releases to AUR, Homebrew, and Chocolatey.
- **Playlist Management**: Add support for creating, saving, and loading playlists.
- **Zoom & Pan**: Fixes for keeping the image centered when zooming in/out are needed.
- **Customisable Key Bindings**: Allow users to customise key bindings for navigation and controls.
- **Image Metadata**: Display image metadata such as resolution, file size, and format.
- **Image Sorting**: Sort images by name, date, or size.
- **Image Filtering**: Filter images by file type or resolution.
- **Plugins**: Add support for Lua plugins to extend the application's functionality.
    - **Package Manager**: Implement a package manager for installing and managing plugins.
- **Thumbnail View**: Display a grid of image thumbnails for quicker navigation.
- **UI Tests**: Write tests to ensure the application functions as expected.

## Contributing

Contributions are welcome! Feel free to submit a pull request to suggest improvements or additional features.

## License

This project is licensed under the MIT License.
