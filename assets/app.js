// Enable an auto-slideshow by adding interval secs to the URL, e.g. `&auto=4`.
window.onhashchange = function() {
  var matches = location.hash.match(new RegExp('auto=([^&]*)'));
  var auto = matches ? matches[1] : null;
  if(window.slideInterval != undefined) {
    clearInterval(window.slideInterval)
  }
  window.slideInterval = setInterval(function(){
    try {
      window.slideshow.next();
    } catch(e) {
      clearInterval(window.slideInterval)
    }
  }, auto*1000);
}

function zoom(imageIndex) {
  var pswpElement = document.querySelectorAll('.pswp')[0];
  var options = {
    index: imageIndex,
    shareEl: false
  };
  window.slideshow = new PhotoSwipe(
    pswpElement,
    PhotoSwipeUI_Default,
    images,
    options
  );
  slideshow.init();
}
