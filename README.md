# miniphoto

A simple and little web photo gallery.

Provides a basic Web UI for browsing photos/images in a given directory. Image
slideshow functionality provided by the [photoswipe][photoswipe] library.

Directory driven, expects the following, simple 2-level directory structure:

```
root
├── some-images
│   ├── credits.html
│   ├── image1.jpg
│   ├── image2.jpg
└── more-images
│   ├── image3.jpg
│   ├── image4.jpg
└── and-so-on
│   ├── image5.jpg
```

If a file named `credits.html` is placed in a photo directory, it's html content
will be rendered at the top of that page.

## Usage

See `bin/dev` & `bin/deploy`

## License

This project uses the MIT License. See [LICENSE](LICENSE).

[photoswipe]: http://photoswipe.com/
