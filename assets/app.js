function zoom(imageIndex) {
  var pswpElement = document.querySelectorAll('.pswp')[0];
  var options = {
    index: imageIndex,
    shareEl: false
  };
  var gallery = new PhotoSwipe(
    pswpElement,
    PhotoSwipeUI_Default,
    images,
    options
  );
  gallery.init();
}
