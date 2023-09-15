$.MyAPP = {};

$.MyAPP.Disconnect = function (cname){
  console.log(cname)
  $.ajax({
    type: "DELETE",
    dataType: "json",
    url: "api/v1/session",
    data: JSON.stringify({ "cname": cname }),
    success: function(data) {
      location.reload();
      console.log(data);
    },
    error: function(a,b,c) {
      console.log(a,b,c)
      location.reload();
    }
  });
}

$.MyAPP.Restart = function (sname){
  console.log(sname)
  $.ajax({
    type: "DELETE",
    dataType: "json",
    url: "/signal",
    data: JSON.stringify({ "sname": sname }),
    success: function(data) {
      location.reload();
      console.log(data);
    },
    error: function(a,b,c) {
      console.log(a,b,c)
      location.reload();
    }
  });
}

$(function() {
  new Clipboard('.button-copy');

  //$( ".btn-disconnect" ).click(function() {
  //  alert( "Handler for .click() called." );
  //});
  //window.location.reload();
})

// Switch theme start

var toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
var currentTheme = localStorage.getItem('theme');
var mainHeader = document.querySelector('.main-header');

if (currentTheme) {
  if (currentTheme === 'dark') {
    if (!document.body.classList.contains('dark-mode')) {
      document.body.classList.add("dark-mode");
    }
    if (mainHeader.classList.contains('navbar-light')) {
      mainHeader.classList.add('navbar-dark');
      mainHeader.classList.remove('navbar-light');
    }
    toggleSwitch.checked = true;
  }
}

function switchTheme(e) {
  if (e.target.checked) {
    if (!document.body.classList.contains('dark-mode')) {
      document.body.classList.add("dark-mode");
    }
    if (mainHeader.classList.contains('navbar-light')) {
      mainHeader.classList.add('navbar-dark');
      mainHeader.classList.remove('navbar-light');
    }
    localStorage.setItem('theme', 'dark');
  } else {
    if (document.body.classList.contains('dark-mode')) {
      document.body.classList.remove("dark-mode");
    }
    if (mainHeader.classList.contains('navbar-dark')) {
      mainHeader.classList.add('navbar-light');
      mainHeader.classList.remove('navbar-dark');
    }
    localStorage.setItem('theme', 'light');
  }
}

toggleSwitch.addEventListener('change', switchTheme, false);

// Switch theme end