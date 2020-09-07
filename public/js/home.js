(function($) {
  "use strict"; // Start of use strict

  // Smooth scrolling using jQuery easing
  $('[href*="#"]:not([href="#"])').click(function() {
    if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
      var target = $(this.hash);
      target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
      if (target.length) {
        $('html, body').animate({
          scrollTop: (target.offset().top - 48)
        }, 1000, "easeInOutExpo");
        return false;
      }
    }
  });

  // Scroll reveal calls
  window.sr = ScrollReveal();
  sr.reveal('[data-sr-scale]', {
    duration: 600,
    scale: 0.3,
    distance: '0px'
  }, 200);
  sr.reveal('[data-sr-slide-up]', {
    distance: '150%',
    origin: 'bottom',
    opacity: 0
  }, 2000);
})(jQuery); // End of use strict
