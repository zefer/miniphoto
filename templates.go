package main

const layoutTemplate = `
{{define "layout"}}
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>{{.Title}}</title>
  <meta name="description" content="">
  <meta name="viewport" content="width=device-width, initial-scale=1">

  <link rel="stylesheet" href="assets/photoswipe.css">
  <link rel="stylesheet" href="assets/default-skin/default-skin.css">
  <link rel="stylesheet" href="assets/app.css">
  <script src="assets/photoswipe.min.js"></script>
  <script src="assets/photoswipe-ui-default.min.js"></script>
  <script src="assets/app.js"></script>

  <link rel="icon" href="assets/icon.png" />
  <link rel="apple-touch-icon-precomposed" href="assets/icon.png">
  <meta name="apple-mobile-web-app-capable" content="yes">
  <meta name="application-name" content="{{.AppName}}">
</head>
<body>

<header>
{{if .Home}}{{else}}<a href="/" title="Home" class="home">&#9664; home</a>{{end}}
<h1>{{.Title}}</h1>
</header>

{{if .Home}}
  <nav><ul>
  {{range .Dirs}}<li><a href="{{ . }}">{{ . }}</a></li>{{else}}<div><strong>no photos yet, check back soon</strong></div>{{end}}
  </ul></nav>
{{else}}
  <main>
	<ul class="gallery">
	{{range $i, $img := .Images}}<li><img src="{{.Src}}" alt="{{.Title}}" height="150" onclick="zoom({{$i}});" /></li>{{end}}
	</ul>
  {{template "photoswipe"}}
  </main>
{{end}}

<script>
var images = {{.ImageJson}};
</script>

</body>
</html>
{{end}}
`

const photoswipeTemplate = `
{{define "photoswipe"}}
<!-- Root element of PhotoSwipe. Must have class pswp. -->
<div class="pswp" tabindex="-1" role="dialog" aria-hidden="true">

  <!-- Background of PhotoSwipe. 
         It's a separate element as animating opacity is faster than rgba(). -->
  <div class="pswp__bg"></div>

  <!-- Slides wrapper with overflow:hidden. -->
  <div class="pswp__scroll-wrap">

    <!-- Container that holds slides. 
            PhotoSwipe keeps only 3 of them in the DOM to save memory.
            Don't modify these 3 pswp__item elements, data is added later on. -->
    <div class="pswp__container">
      <div class="pswp__item"></div>
      <div class="pswp__item"></div>
      <div class="pswp__item"></div>
    </div>

    <!-- Default (PhotoSwipeUI_Default) interface on top of sliding area. Can be changed. -->
    <div class="pswp__ui pswp__ui--hidden">

      <div class="pswp__top-bar">

        <!--  Controls are self-explanatory. Order can be changed. -->

        <div class="pswp__counter"></div>

        <button class="pswp__button pswp__button--close" title="Close (Esc)"></button>

        <button class="pswp__button pswp__button--share" title="Share"></button>

        <button class="pswp__button pswp__button--fs" title="Toggle fullscreen"></button>

        <button class="pswp__button pswp__button--zoom" title="Zoom in/out"></button>

        <!-- Preloader demo http://codepen.io/dimsemenov/pen/yyBWoR -->
        <!-- element will get class pswp__preloader--active when preloader is running -->
        <div class="pswp__preloader">
          <div class="pswp__preloader__icn">
            <div class="pswp__preloader__cut">
              <div class="pswp__preloader__donut"></div>
            </div>
          </div>
        </div>
      </div>

      <div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
        <div class="pswp__share-tooltip"></div> 
      </div>

      <button class="pswp__button pswp__button--arrow--left" title="Previous (arrow left)">
      </button>

      <button class="pswp__button pswp__button--arrow--right" title="Next (arrow right)">
      </button>

      <div class="pswp__caption">
        <div class="pswp__caption__center"></div>
      </div>

    </div>

  </div>

</div>
{{end}}
`
