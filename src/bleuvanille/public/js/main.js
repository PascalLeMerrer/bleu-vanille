var contactEmail = "info" + "@" + "bleuvanille.fr"
var startTime;

/* Hides alerts on page loading and init scrolling */
$(document).ready(function() {
  $('#success-alert').hide();
  $('#error-alert').hide();
  initScrollspy();
  startTime = new Date();
  addScroll('#menuButton')
  addScroll('#arrowToMenu')
  $('#contact').text(contactEmail);
});

/* Adds scrolling to the first section when clicking on a given link */
function addScroll(id) {
  $(id).click(function(e){
      var topMenu = $("#top-menu"),
          offsetTop = $('#menu').offset().top - topMenu.outerHeight() + 1;
      $('html, body').stop().animate({
        scrollTop: offsetTop
      }, 400);
      e.preventDefault();
  });
}

/* Highlights the current menu item and add scrolling to the corresponding section
 */

function initScrollspy() {

  // Cache selectors
  var lastId,
    topMenu = $("#top-menu"),
    topMenuHeight = topMenu.outerHeight(),
    // All list items
    menuItems = topMenu.find("a"),
    // Anchors corresponding to menu items
    scrollItems = menuItems.map(function() {
      var item = $($(this).attr("href"));
      if (item.length) {
        return item;
      }
    });

  // Bind click handler to menu items
  // so we can get a fancy scroll animation
  menuItems.click(function(e) {
    var href = $(this).attr("href"),
      offsetTop = href === "#" ? 0 : $(href).offset().top - topMenuHeight + 1;
    $('html, body').stop().animate({
      scrollTop: offsetTop
    }, 400);
    e.preventDefault();
  });

  // Bind to scroll
  $(window).scroll(function() {
    // Get container scroll position
    var fromTop = $(this).scrollTop() + topMenuHeight;

    // Get id of current scroll item
    var cur = scrollItems.map(function() {
      if ($(this).offset().top < fromTop)
        return this;
    });
    // Get the id of the current element
    cur = cur[cur.length - 1];
    var id = cur && cur.length ? cur[0].id : "";
    if (lastId !== id) {
      $('#' + lastId + "Link").removeClass("menu-selected");
      $('#' + id + "Link").addClass("menu-selected");
      lastId = id;
    }
  });
}
/* -------------------------------------------------------------------------------------------*/
/* Registers the email address of the user if it is valid */
$("#submit").click(function(event) {
  var regex = /^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/;
  var email = $('#emailInput').val();
  var timeSpentSeconds = Math.round((new Date() - startTime)/1000); // time spent on the page
  if (regex.test(email)) {
    var body = {
      email: email,
      timeSpent: timeSpentSeconds
    }
    var posting = $.post('/contacts', body);
    posting.done(function(response) {
      $('#success-alert').text("Votre adresse email a été enregistrée. Nous vous donnerons des nouvelles de Bleu Vanille prochainement.")
      $('#success-alert').show().delay(2000).fadeOut(1000);
    });
    posting.fail(function(response) {
      if (response.status == 409) {
        $('#error-alert').text("Vous êtes déjà inscrit.")
      } else {
        $('#error-alert').text("Une erreur s'est produite. Veuillez réessayer dans quelques minutes. Si l'erreur persiste, prévenez-nous à " + contactEmail)
      }
      $('#error-alert').show().delay(2000).fadeOut(1000);
    });
  } else {
    $('#error-alert').text("Cela ne ressemble pas à une adresse email valide.")
    $('#error-alert').show().delay(2000).fadeOut(1000);
  }
  event.preventDefault();
});
/* -------------------------------------------------------------------------------------------*/
/* Change the user password if it is valid */
$("#passwordResetForm").submit(function(event) {
  var password = $('#passwordInput').val();
  if (password && password.length > 8) {
    header = "Bearer " + $('#token').val()
    var posting = $.ajax({
        url: '/special/resetPassword',
        data: { password: password },
        beforeSend: function(xhr){xhr.setRequestHeader('Authorization', header);},
        type: "POST"
    });

    posting.done(function(response) {
      console.log("Votre mot de passe a été modifié. Vous pouvez l'utiliser pour vous connecter.");
      $('#success-alert').text("Votre mot de passe a été modifié. Vous pouvez l'utiliser pour vous connecter.");
      $('#success-alert').show().delay(2000).fadeOut(1000);
    });
    posting.fail(function(response) {
      console.log("Une erreur s'est produite. Veuillez réessayer dans quelques minutes.");
      $('#error-alert').text("Une erreur s'est produite. Veuillez réessayer dans quelques minutes.");
      $('#error-alert').show().delay(2000).fadeOut(1000);
    });
  } else {
    console.log("Votre mot de passe est trop court.");
    $('#error-alert').text("Votre mot de passe est trop court.");
    $('#error-alert').show().delay(2000).fadeOut(1000);
  }
  event.preventDefault();
});
